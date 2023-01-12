package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"

	_ "embed"

	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/dns"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"

	"github.com/ojo-network/ojo/infra/pulumi/testnet/unit"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"github.com/umee-network/umee-infra/infra/pulumi/common/components/caddy"
	"github.com/umee-network/umee-infra/infra/pulumi/common/components/cosmosuptime"
	"github.com/umee-network/umee-infra/infra/pulumi/common/resources"
	"github.com/umee-network/umee-infra/lib/caddyconfiggen"
	"github.com/umee-network/umee-infra/lib/umeedconfiggen"
	netconfig "github.com/umee-network/umee-infra/lib/umeednetworkconfigurator"
	"github.com/umee-network/umee-infra/lib/umeedwrapper"
)

func (network Network) Provision(ctx *pulumi.Context, secrets []NodeSecretConfig) error {
	// addrs
	var addrs pulumi.StringArray
	var nodeHostNames []string
	var dnsArgs []DNSRecordsArgs

	if len(network.NodeConfig.Locations) == 0 {
		log.Fatal("must specify at least one location in nodeConfig")
	}

	for i := 0; i < network.NumNodes; i++ {
		moniker := genMoniker(network.ChainID, i)
		location := pickLocation(network.NodeConfig.Locations, i)
		addr, err := compute.NewAddress(ctx, moniker+"-ip", &compute.AddressArgs{
			Labels: pulumi.StringMap{
				"chain_id": pulumi.String(network.ChainID),
			},
			NetworkTier: pulumi.String("STANDARD"),
			Region:      pulumi.String(location.Region),
		})
		if err != nil {
			return err
		}

		addrs = append(addrs, addr.Address)

		nodeHostName := fmt.Sprintf("%s.%s.node.ojo.network", moniker, network.ChainID)
		nodeHostNames = append(nodeHostNames, nodeHostName)

		dnsArg := DNSRecordsArgs{
			APIHostName:  fmt.Sprintf("%s.%s", "api", nodeHostName),
			RPCHostName:  fmt.Sprintf("%s.%s", "rpc", nodeHostName),
			GRPCHostName: fmt.Sprintf("%s.%s", "grpc", nodeHostName),
		}
		dnsArgs = append(dnsArgs, dnsArg)

		_, err = createDNSRecords(ctx, addr.Address, dnsArg)
		if err != nil {
			return err
		}
	}

	netPackResult, err := network.performGenesisNetpack(
		ctx,
		addrs,
		network,
	)
	if err != nil {
		return err
	}
	ctx.Export("netpack-result", netPackResult)

	for i := 0; i < network.NumNodes; i++ {
		moniker := genMoniker(network.ChainID, i)
		location := pickLocation(network.NodeConfig.Locations, i)
		bootDisk := &compute.InstanceBootDiskArgs{
			DeviceName: pulumi.String(fmt.Sprintf("%s-bootdisk", moniker)),
			InitializeParams: &compute.InstanceBootDiskInitializeParamsArgs{
				Image: pulumi.String("family/ubuntu-minimal-2204-lts"),
				Type:  pulumi.String(network.NodeConfig.DiskType),
				Size:  pulumi.Int(network.NodeConfig.DiskSizeGB),
			},
		}

		conf := config.New(ctx, "")
		sshPublic := conf.Require("sshpublic")
		sshPrivate := conf.RequireSecret("sshprivate").ApplyT(func(b64private string) (string, error) {
			privatebytes, err := base64.StdEncoding.DecodeString(b64private)
			if err != nil {
				return "", err
			}

			return string(privatebytes), nil
		}).(pulumi.StringOutput)

		ubuntuPubkey := pulumi.String("ubuntu:" + sshPublic)

		serviceAccount, err := createServiceAccount(ctx, moniker+"-svc", "service account for "+moniker)
		if err != nil {
			return err
		}

		startupScript := pulumi.String(genStartupScript())
		instance, err := compute.NewInstance(ctx, moniker+"-instance", &compute.InstanceArgs{
			Name: pulumi.String(moniker),
			Labels: pulumi.StringMap{
				"chain_id": pulumi.String(network.ChainID),
			},
			MachineType:            pulumi.String(network.NodeConfig.MachineType),
			Zone:                   pulumi.String(location.Zone),
			Hostname:               pulumi.String(nodeHostNames[i]),
			AllowStoppingForUpdate: pulumi.Bool(true),
			BootDisk:               bootDisk,
			MetadataStartupScript:  startupScript,
			Metadata: pulumi.StringMap{
				"ssh-keys":               ubuntuPubkey,
				"block-project-ssh-keys": pulumi.String("true"),
			},
			ServiceAccount: &compute.InstanceServiceAccountArgs{
				Email: serviceAccount.Email,
				Scopes: pulumi.StringArray{
					pulumi.String("cloud-platform"),
				},
			},
			NetworkInterfaces: compute.InstanceNetworkInterfaceArray{
				&compute.InstanceNetworkInterfaceArgs{
					Network: pulumi.String("default"),
					AccessConfigs: compute.InstanceNetworkInterfaceAccessConfigArray{
						compute.InstanceNetworkInterfaceAccessConfigArgs{
							NatIp:       addrs[i],
							NetworkTier: pulumi.String("STANDARD"),
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}

		conn := remote.ConnectionArgs{
			Host:       addrs.ToStringArrayOutput().Index(pulumi.Int(i)),
			Port:       pulumi.Float64(22),
			User:       pulumi.String("ubuntu"),
			PrivateKey: sshPrivate,
		}

		startupScriptsComplete, err := remote.NewCommand(
			ctx,
			moniker+"-bootstrap-script-wait-until-ready",
			&remote.CommandArgs{
				Triggers:   pulumi.Array{startupScript},
				Connection: conn,
				Update:     pulumi.String("echo updates disabled..."),
				Create: pulumi.Sprintf(`
                  for VARIABLE in 1 2 3 4 5 6 7 8 9 .. N
                  do
                    if test -f "/tmp/STARTUP_FINISHED"; then
                      exit 0
                    else
                      echo 'System startup script incomplete; sleeping 30 seconds...'
                      sleep 45
                    fi
                  done

                  echo 'Machine is not ready or system startup script did not complete (timeout)'
                  exit 1
                `),
			},
			pulumi.DependsOn([]pulumi.Resource{instance}),
			pulumi.Timeouts(&pulumi.CustomTimeouts{Create: "10m"}),
		)
		if err != nil {
			return err
		}

		techName := network.CosmosHomeFolderName[1:]
		unitSpec := unit.UnitSpec{
			Name:              techName,
			Description:       fmt.Sprintf("%s daemon", techName),
			User:              "blockchain",
			BinaryInstallPath: fmt.Sprintf("/usr/local/bin/%s", network.LocalCosmosBinaryPath),
		}
		unit := unitSpec.ToUnit()

		nodeIndexClosure := i

		// custom config functionality before starting unit
		var configReady pulumi.Resource
		tarBallCopyReady, err := remote.NewCopyFile(ctx, moniker+"-"+unit.Name+"-cp-home-tarball", &remote.CopyFileArgs{
			Connection: conn,
			LocalPath: netPackResult.ApplyT(func(val interface{}) string {
				netPackResult := val.(*netconfig.PackResult)
				return netPackResult.Nodes[nodeIndexClosure].TarballPath
			}).(pulumi.StringOutput),
			RemotePath: pulumi.String("/home/ubuntu/cosmos-home.tgz"),
		}, pulumi.DependsOn([]pulumi.Resource{startupScriptsComplete}))
		if err != nil {
			return err
		}

		homeFolder := "/home/blockchain/" + network.CosmosHomeFolderName
		configReady, err = remote.NewCommand(
			ctx,
			moniker+"-"+unit.Name+"-untar-home",
			&remote.CommandArgs{
				Connection: conn,
				Create: pulumi.Sprintf(`
						    set -e
							sudo adduser -q --disabled-password --gecos "Blockchain Non-privileged User" blockchain
							sudo -u blockchain mkdir -p %s
							sudo cp /home/ubuntu/cosmos-home.tgz /home/blockchain/
							sudo -u blockchain tar -xvzf /home/blockchain/cosmos-home.tgz -C %s
						`, homeFolder, homeFolder),
			}, pulumi.DependsOn([]pulumi.Resource{tarBallCopyReady}),
		)
		if err != nil {
			return err
		}

		uploadCosmosBinary, err := remote.NewCopyFile(ctx, moniker+"-"+unit.Name+"-cp-cosmos-binary", &remote.CopyFileArgs{
			Connection: conn,
			// TODO: don't assume /usr/local/ as the base path (brittle); will work for now since we control action file, may not work on a particular devs machine
			LocalPath:  pulumi.Sprintf("/usr/local/bin/%s", network.LocalCosmosBinaryPath),
			RemotePath: pulumi.Sprintf("/home/ubuntu/%s", network.LocalCosmosBinaryPath),
		}, pulumi.DependsOn([]pulumi.Resource{configReady}))
		if err != nil {
			return err
		}

		installCosmosBinary, err := remote.NewCommand(
			ctx,
			moniker+"-"+unit.Name+"-install-cosmos-binary",
			&remote.CommandArgs{
				Connection: conn,
				Create: pulumi.Sprintf(`
						    set -e
							sudo cp /home/ubuntu/%s /usr/local/bin/
							sudo chmod a+x /usr/local/bin/%s
						`, network.LocalCosmosBinaryPath, network.LocalCosmosBinaryPath),
			}, pulumi.DependsOn([]pulumi.Resource{uploadCosmosBinary}),
		)
		if err != nil {
			return err
		}

		unitBody := unit.GenSystemdUnit()
		unitPath := pulumi.String(path.Join("/etc/systemd/system", unit.Name+".service"))
		unitInstall, err := resources.NewStringToRemoteFileCommand(ctx, moniker+"-"+unit.Name+"-systemd-unit", resources.StringToRemoteFileCommandArgs{
			Connection:      conn,
			Body:            unitBody,
			DestinationPath: unitPath,
			FileMode:        pulumi.String("0644"),
			FileUser:        unit.User,
			FileGroup:       unit.User,
			FolderMode:      pulumi.String("0755"),
			FolderUser:      unit.User,
			FolderGroup:     unit.User,
			RunAfter:        pulumi.Sprintf("sudo systemctl daemon-reload && sudo systemctl enable %s", unit.Name),
			Triggers:        pulumi.Array{unitPath, unitBody},
		}, pulumi.DependsOn([]pulumi.Resource{configReady, installCosmosBinary}))
		if err != nil {
			return err
		}

		caddy, err := createCaddy(ctx, conn, moniker, dnsArgs[i], []pulumi.Resource{startupScriptsComplete}, pulumi.Array{startupScript})
		if err != nil {
			return err
		}

		uptimeMonitoring, err := network.uptimeMonitoring(ctx, moniker, dnsArgs[i], []pulumi.Resource{})
		if err != nil {
			return err
		}

		rebootDeps := []pulumi.Resource{
			caddy,
			configReady,
			unitInstall,
			uptimeMonitoring,
			installCosmosBinary,
		}

		_, err = remote.NewCommand(
			ctx,
			moniker+"-reboot",
			&remote.CommandArgs{
				Connection: conn,
				Update:     pulumi.String("echo updates disabled..."),
				Create:     pulumi.String("sleep 30 && sudo shutdown -r 1"),
			},
			pulumi.DependsOn(rebootDeps),
		)
		if err != nil {
			return err
		}
	}

	ctx.Export("node-hostnames", pulumi.ToStringArray(nodeHostNames))

	return nil
}

func (n Network) performGenesisNetpack(ctx *pulumi.Context, addrs pulumi.StringArray, network Network) (pulumi.Output, error) {
	stackName := strings.Join([]string{"ojo-network", ctx.Project(), ctx.Stack()}, "/")
	stack, err := pulumi.NewStackReference(ctx, stackName, nil)
	if err != nil {
		return nil, err
	}

	return pulumi.All(
		addrs.ToStringArrayOutput(),
		stack.GetOutput(pulumi.String("netpack-result")),
	).ApplyT(func(args []interface{}) (*netconfig.PackResult, error) {
		addrs := args[0].([]string)
		_netpack, ok := args[1].(map[string]interface{})

		if ok {
			__netpack, err := json.Marshal(_netpack)
			if err != nil {
				return nil, err
			}
			var netpack netconfig.PackResult
			err = json.Unmarshal(__netpack, &netpack)
			if err != nil {
				return nil, err
			}
			log.Println("using existing genesis netpack...")
			return &netpack, nil
		}

		netConfig := netconfig.NetworkConfig{
			GeneratePersistentPeers: true,
			ChainID:                 network.ChainID,
			NumNodes:                network.NumNodes,
			// implement in yaml only support primitive ; don't support complex types
			GenesisModifierFunc: func(_ *netconfig.Network, genesis string) (string, error) {
				return n.NetworkGenesisMutations.MutateGenesis(genesis)
			},
			// TODO: may need to allow more customization per blockchain on config but avoid for now
			Configs: func() []umeedconfiggen.Config {
				var out []umeedconfiggen.Config
				for i, addr := range addrs {
					config := umeedconfiggen.
						NewDefaultConfig().
						SetExternalAddress(addr + ":26656").
						SetMoniker(genMoniker(network.ChainID, i))

					out = append(out, config)
				}

				return out
			}(),
			// TODO: may need to allow more customization per blockchain on config but avoid for now
			AppConfigs: func() []umeedconfiggen.AppConfig {
				var out []umeedconfiggen.AppConfig
				for range addrs {
					config := umeedconfiggen.NewDefaultAppConfig()
					out = append(out, config)
				}

				return out
			}(),
			NodeGenesisAccounts: network.NodeGenesisAccounts,
			GenesisAccounts:     network.GenesisAccounts,
		}
		netWrapper := umeedwrapper.Wrapper{
			CosmosBinaryPath: network.LocalCosmosBinaryPath,
		}
		netconfigurator, err := netconfig.NewNetwork(netConfig, netWrapper)
		if err != nil {
			return nil, err
		}
		netPackResult, err := netconfigurator.ConfigureAndPack("/tmp")
		if err != nil {
			return nil, err
		}

		return &netPackResult, nil
	}), nil
}

func (n Network) uptimeMonitoring(
	ctx *pulumi.Context,
	moniker string,
	dns DNSRecordsArgs,
	dependsOn []pulumi.Resource,
) (*cosmosuptime.CosmosUptime, error) {
	return cosmosuptime.NewCosmosUptime(ctx, moniker+"-mon", cosmosuptime.CosmosUptimeArgs{
		APIHostname:  pulumi.String(dns.APIHostName),
		RPCHostname:  pulumi.String(dns.RPCHostName),
		GRPCHostname: pulumi.String(dns.GRPCHostName),
	}, pulumi.DependsOn(dependsOn))
}

func createServiceAccount(ctx *pulumi.Context, name string, desc string) (*serviceaccount.Account, error) {
	account, err := serviceaccount.NewAccount(ctx, name, &serviceaccount.AccountArgs{
		AccountId:   pulumi.String(name),
		DisplayName: pulumi.String(name),
		Description: pulumi.String(desc),
	})
	if err != nil {
		return nil, err
	}

	iamMember := account.Email.ApplyT(func(email string) string {
		return "serviceAccount:" + email
	}).(pulumi.StringOutput)

	gcpProject, ok := ctx.GetConfig("gcp:project")
	if !ok {
		return nil, fmt.Errorf("gcp:project must be set")
	}

	_, err = projects.NewIAMMember(ctx, name+"-metricwriter-role", &projects.IAMMemberArgs{
		Role:    pulumi.String("roles/monitoring.metricWriter"),
		Member:  iamMember,
		Project: pulumi.String(gcpProject),
	})
	if err != nil {
		return nil, err
	}

	_, err = projects.NewIAMMember(ctx, name+"-logwriter-role", &projects.IAMMemberArgs{
		Role:    pulumi.String("roles/logging.logWriter"),
		Member:  iamMember,
		Project: pulumi.String(gcpProject),
	})
	if err != nil {
		return nil, err
	}

	return account, nil
}

func createCaddy(
	ctx *pulumi.Context,
	conn remote.ConnectionArgs,
	moniker string,
	dns DNSRecordsArgs,
	dependsOn []pulumi.Resource,
	recreateTriggers pulumi.ArrayInput,
) (*caddy.CaddyDaemon, error) {
	return caddy.NewCaddyDaemon(ctx, moniker+"-caddy", caddy.CaddyDaemonArgs{
		Connection: conn,
		CaddyConfig: caddyconfiggen.Config{
			LocalProxyApps: []caddyconfiggen.LocalProxyApp{
				{
					DomainName: dns.APIHostName,
					LocalPort:  1317,
				},
				{
					DomainName: dns.RPCHostName,
					LocalPort:  26657,
				},
				{
					DomainName: dns.GRPCHostName,
					LocalPort:  9090,
					IsGRPC:     true,
				},
			},
		},
		Triggers: recreateTriggers,
	}, pulumi.DependsOn(dependsOn))
}

type DNSRecordsArgs struct {
	APIHostName  string
	RPCHostName  string
	GRPCHostName string
}

func createDNSRecords(ctx *pulumi.Context, ip pulumi.StringInput, names DNSRecordsArgs) ([]pulumi.Resource, error) {
	project := "ojo-network"
	managedZone := "ojo-network"

	apiDNS, err := dns.NewRecordSet(ctx, "DNS-A-"+names.APIHostName, &dns.RecordSetArgs{
		ManagedZone: pulumi.String(managedZone),
		Project:     pulumi.String(project),
		Type:        pulumi.String("A"),
		Name:        pulumi.String(names.APIHostName + "."),
		Rrdatas:     pulumi.StringArray{ip},
		Ttl:         pulumi.Int(300),
	})
	if err != nil {
		return nil, err
	}

	rpcDNS, err := dns.NewRecordSet(ctx, "DNSA-"+names.RPCHostName, &dns.RecordSetArgs{
		ManagedZone: pulumi.String(managedZone),
		Project:     pulumi.String(project),
		Type:        pulumi.String("A"),
		Name:        pulumi.String(names.RPCHostName + "."),
		Rrdatas:     pulumi.StringArray{ip},
		Ttl:         pulumi.Int(300),
	})
	if err != nil {
		return nil, err
	}

	grpcDNS, err := dns.NewRecordSet(ctx, "DNSA-"+names.GRPCHostName, &dns.RecordSetArgs{
		ManagedZone: pulumi.String(managedZone),
		Project:     pulumi.String(project),
		Type:        pulumi.String("A"),
		Name:        pulumi.String(names.GRPCHostName + "."),
		Rrdatas:     pulumi.StringArray{ip},
		Ttl:         pulumi.Int(300),
	})
	if err != nil {
		return nil, err
	}

	return []pulumi.Resource{apiDNS, rpcDNS, grpcDNS}, nil
}

func genMoniker(chainID string, index int) string {
	return fmt.Sprintf("devnet-n%d", index)
}

func pickLocation(locations []NodeLocation, index int) NodeLocation {
	return locations[index%len(locations)]
}
