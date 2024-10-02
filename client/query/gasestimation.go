package query

import (
	"context"

	gasestimatetypes "github.com/ojo-network/ojo/x/gasestimate/types"
)

// GovQueryClient returns the govtypes.QueryClient
// initialized with the clients grpc connection
func (c *Client) GasEstimateClient() gasestimatetypes.QueryClient {
	return gasestimatetypes.NewQueryClient(c.grpcConn)
}

func (c *Client) GetGasEstimate(ctx context.Context, network string) (int64, error) {
	queryClient := c.GasEstimateClient()
	query := &gasestimatetypes.GasEstimateRequest{
		Network: network,
	}

	res, err := queryClient.GasEstimate(ctx, query)
	if err != nil {
		return 0, err
	}

	return res.GasEstimate, nil
}
