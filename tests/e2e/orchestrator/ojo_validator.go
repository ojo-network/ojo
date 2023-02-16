package orchestrator

import (
	"context"
	"fmt"

	"github.com/ory/dockertest/v3"

	appparams "github.com/ojo-network/ojo/app/params"
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
	name           string
	mnemonic       string
	chainID        string
	dockerResource *dockertest.Resource
	rpc            *rpchttp.HTTP
}

func (o *Orchestrator) initOjoVal(index int) (val *Validator, err error) {
	val = &Validator{
		name:    fmt.Sprintf("ojoVal%d", index),
		chainID: ojo_chain_id,
	}
	val.mnemonic, err = createMnemonic()
	return
}

func (o *Orchestrator) startOjoVal(val *Validator, configDir string) (err error) {
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

	return val.setTendermintEndpoint()
}

func (val *Validator) setTendermintEndpoint() (err error) {
	path := val.dockerResource.GetHostPort(fmt.Sprintf("%s/tcp", ojo_tmrpc_port))
	endpoint := fmt.Sprintf("tcp://%s", path)
	val.rpc, err = rpchttp.New(endpoint, "/websocket")
	return
}

func (val *Validator) BlockHeight() (int64, error) {
	status, err := val.rpc.Status(context.Background())
	if err != nil {
		return 0, err
	}
	return status.SyncInfo.LatestBlockHeight, nil
}
