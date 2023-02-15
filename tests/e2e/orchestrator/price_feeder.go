package orchestrator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ory/dockertest/v3"
)

const (
	price_feeder_container_name = "price-feeder"
	price_feeder_container_repo = "ghcr.io/ojo-network"
	price_feeder_server_port    = "7171/tcp"
)

type PriceFeeder struct {
	dockerResource *dockertest.Resource
	endpoint       string
}

func (o *Orchestrator) CreatePriceFeeder(validator *Validator) (pf *PriceFeeder, err error) {
	pf.dockerResource, err = o.dockerPool.RunWithOptions(
		&dockertest.RunOptions{
			Name:       price_feeder_container_name,
			Repository: price_feeder_container_repo,
			NetworkID:  o.dockerNetwork.Network.ID,
			Env: []string{
				fmt.Sprintf("GRPC_ENDPOINT=%s", validator.chainID),
				fmt.Sprintf("TMRPC_ENDPOINT=%s", validator.rpc),
			},
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

	pf.setEndpoint()
	return
}

func (pf *PriceFeeder) setEndpoint() {
	pf.endpoint = fmt.Sprintf("http://%s", pf.dockerResource.GetHostPort(price_feeder_server_port))
}

func (pf *PriceFeeder) GetPrices() (prices map[string]interface{}, err error) {
	url := fmt.Sprintf("%s/api/v1/prices", pf.endpoint)

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var respBody map[string]interface{}
	if err = json.Unmarshal(bz, &respBody); err != nil {
		return
	}

	prices, ok := respBody["prices"].(map[string]interface{})
	if !ok {
		err = fmt.Errorf("price feeder: no prices")
	}
	return
}
