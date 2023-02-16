package orchestrator

import (
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
)

// Orchestrator is responsible for managing docker resources,
// their configuration files, and environment variables.
type Orchestrator struct {
	dockerPool    *dockertest.Pool
	dockerNetwork *dockertest.Network

	validators  []*Validator
	PriceFeeder *PriceFeeder
}

func (o *Orchestrator) InitDockerResources(t *testing.T) error {
	var err error

	t.Log("-> initializing docker network")
	err = o.initNetwork()
	if err != nil {
		return err
	}

	t.Log("-> initializing ojo validator")
	new_validator, err := o.CreateOjoValidator(0)
	if err != nil {
		return err
	}
	o.validators = append(o.validators, new_validator)

	t.Log("-> verifying ojo node is creating blocks")
	require.Eventually(
		t,
		func() bool {
			blockHeight, err := new_validator.BlockHeight()
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
	//o.NewPriceFeeder()

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
