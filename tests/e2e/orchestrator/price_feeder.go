package orchestrator

import (
	"fmt"

	"github.com/ory/dockertest/v3"
)

// Creates a new dockertest resource for the price feeder
// and verifies that it is running.
func (o *Orchestrator) initPriceFeeder() {
	var err error

	o.priceFeederResource, err = o.dockerPool.RunWithOptions(
		&dockertest.RunOptions{
			Name:       price_feeder_container_name,
			Repository: price_feeder_container_name,
			NetworkID:  o.dockerNetwork.Network.ID,
			Env: []string{
				fmt.Sprintf("GRPC_ENDPOINT=%s", o.ojoChain.chainId),
				fmt.Sprintf("TMRPC_ENDPOINT=%s", o.priceFeeder.mnemonic),
			},
			Entrypoint: []string{
				"sh",
				"-c",
				"chmod +x .ojo/config/ojo_bootstrap.sh && .ojo/config/ojo_bootstrap.sh",
			},
		},
		noRestart,
	)
	o.Require().NoError(err)

	err = o.setPriceFeederEndpoint()
	o.Require().NoError(err)
}
