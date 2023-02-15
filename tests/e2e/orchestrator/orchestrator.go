package orchestrator

import (
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"

	appparams "github.com/ojo-network/ojo/app/params"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

const (
	ojo_container_name = "ojo"
	ojo_tmrpc_port     = "26657"
	ojo_grpc_port      = "9090"

	price_feeder_container_name = "price-feeder"
	price_feeder_server_port    = "8080"
)

var (
	minGasPrice = appparams.ProtocolMinGasPrice.String()
)

// Orchestrator is responsible for managing docker resources,
// their configuration files, and environment variables.
type Orchestrator struct {
	dockerPool    *dockertest.Pool
	dockerNetwork *dockertest.Network

	ojoResource *dockertest.Resource
	ojoRPC      *rpchttp.HTTP
	ojoChain    *Chain

	priceFeederResource *dockertest.Resource
}

func (o *Orchestrator) InitDockerResources(t *testing.T) error {
	var err error

	t.Log("-> initializing docker network")
	err = o.initNetwork()
	if err != nil {
		return err
	}

	t.Log("-> initializing ojo validator")
	o.initOjod()

	t.Log("-> verifying ojo node is creating blocks")
	require.Eventually(
		t,
		func() bool {
			blockHeight, err := o.ojoBlockHeight()
			if err != nil {
				return false
			}
			return blockHeight >= 3
		},
		time.Minute,
		time.Second*2,
		"ojo node failed to produce blocks",
	)

	t.Log("-> initializing price-feeder")

	return nil
}

func (o *Orchestrator) TearDownDockerResources() error {
	return o.dockerPool.Client.RemoveNetwork(o.dockerNetwork.Network.ID)
}

func (o *Orchestrator) initNetwork() error {
	var err error
	o.dockerPool, err = dockertest.NewPool("")
	if err != nil {
		return err
	}

	o.dockerNetwork, err = o.dockerPool.CreateNetwork("e2e_test_network")
	if err != nil {
		return err
	}
	return nil
}

func noRestart(config *docker.HostConfig) {
	config.RestartPolicy = docker.RestartPolicy{
		Name: "no",
	}
}
