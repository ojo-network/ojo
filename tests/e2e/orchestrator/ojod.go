package orchestrator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ory/dockertest/v3"

	tmconfig "github.com/tendermint/tendermint/config"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

func (o *Orchestrator) initOjod() error {
	var err error

	configDir, err := o.initOjoConfigs()
	if err != nil {
		return err
	}

	o.ojoChain = NewChain("test-ojo")

	o.ojoResource, err = o.dockerPool.RunWithOptions(
		&dockertest.RunOptions{
			Name:       ojo_container_name,
			Repository: ojo_container_name,
			NetworkID:  o.dockerNetwork.Network.ID,
			Env: []string{
				fmt.Sprintf("OJO_CHAIN_ID=%s", o.ojoChain.chainId),
				fmt.Sprintf("MNEMONIC=%s", o.ojoChain.mnemonic),
			},
			Mounts: []string{fmt.Sprintf("%s/:/app/.ojod/config", configDir)},
			Entrypoint: []string{
				"sh",
				"-c",
				"chmod +x .ojod/config/ojo_bootstrap.sh && .ojod/config/ojo_bootstrap.sh",
			},
		},
		noRestart,
	)
	if err != nil {
		return err
	}

	err = o.setTendermintEndpoint()
	return err
}

func (o *Orchestrator) initOjoConfigs() (dir string, err error) {
	dir, err = os.MkdirTemp("", "e2e-configs")
	if err != nil {
		return
	}

	_, err = copyFile(
		filepath.Join("./config/", "ojo_bootstrap.sh"),
		filepath.Join(dir, "ojo_bootstrap.sh"),
	)
	if err != nil {
		return
	}

	configPath := filepath.Join(dir, "config.toml")
	config := tmconfig.DefaultConfig()
	config.P2P.ListenAddress = "tcp://0.0.0.0:26656"
	config.RPC.ListenAddress = fmt.Sprintf("tcp://0.0.0.0:%s", ojo_tmrpc_port)
	config.StateSync.Enable = false
	config.P2P.AddrBookStrict = false
	config.P2P.Seeds = ""
	tmconfig.WriteConfigFile(configPath, config)

	return
}

func (o *Orchestrator) ojoBlockHeight() (int64, error) {
	status, err := o.ojoRPC.Status(context.Background())
	if err != nil {
		return 0, err
	}
	return status.SyncInfo.LatestBlockHeight, nil
}

func (o *Orchestrator) setTendermintEndpoint() (err error) {
	path := o.ojoResource.GetHostPort(fmt.Sprintf("%s/tcp", ojo_tmrpc_port))
	endpoint := fmt.Sprintf("tcp://%s", path)
	o.ojoRPC, err = rpchttp.New(endpoint, "/websocket")
	return
}
