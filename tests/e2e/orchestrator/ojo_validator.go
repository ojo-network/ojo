package orchestrator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ory/dockertest/v3"

	srvconfig "github.com/cosmos/cosmos-sdk/server/config"
	appparams "github.com/ojo-network/ojo/app/params"
	tmconfig "github.com/tendermint/tendermint/config"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

const (
	ojo_container_name = "ojo"
	ojo_tmrpc_port     = "26657"
	ojo_grpc_port      = "9090"
	ojo_chain_id       = "ojo-test-1"
)

var (
	ojoMinGasPrice = appparams.ProtocolMinGasPrice.String()
)

type Validator struct {
	mnemonic       string
	chainID        string
	dockerResource *dockertest.Resource
	rpc            *rpchttp.HTTP
}

func (o *Orchestrator) CreateOjoValidator() (val *Validator, err error) {
	val = &Validator{}
	val.chainID = ojo_chain_id
	val.mnemonic, err = createMnemonic()
	if err != nil {
		return
	}

	configDir, err := initOjoConfigs()
	if err != nil {
		return
	}

	val.dockerResource, err = o.dockerPool.RunWithOptions(
		&dockertest.RunOptions{
			Name:       ojo_container_name,
			Repository: ojo_container_name,
			NetworkID:  o.dockerNetwork.Network.ID,
			Env: []string{
				fmt.Sprintf("OJO_CHAIN_ID=%s", val.chainID),
				fmt.Sprintf("MNEMONIC=%s", val.mnemonic),
			},
			Mounts: []string{fmt.Sprintf("%s/:/ojo/.ojo/config", configDir)},
			Entrypoint: []string{
				"sh",
				"-c",
				"chmod +x .ojo/config/ojo_bootstrap.sh && .ojo/config/ojo_bootstrap.sh",
			},
		},
		noRestart,
	)
	if err != nil {
		return
	}

	err = val.setTendermintEndpoint()
	return
}

func initOjoConfigs() (dir string, err error) {
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

	appCfgPath := filepath.Join(dir, "app.toml")
	appConfig := srvconfig.DefaultConfig()
	appConfig.API.Enable = true
	appConfig.MinGasPrices = ojoMinGasPrice
	srvconfig.WriteConfigFile(appCfgPath, appConfig)

	return
}

func (val *Validator) BlockHeight() (int64, error) {
	status, err := val.rpc.Status(context.Background())
	if err != nil {
		return 0, err
	}
	return status.SyncInfo.LatestBlockHeight, nil
}

func (val *Validator) setTendermintEndpoint() (err error) {
	path := val.dockerResource.GetHostPort(fmt.Sprintf("%s/tcp", ojo_tmrpc_port))
	endpoint := fmt.Sprintf("tcp://%s", path)
	val.rpc, err = rpchttp.New(endpoint, "/websocket")
	return
}
