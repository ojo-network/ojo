package query

import (
	"context"

	airdroptypes "github.com/ojo-network/ojo/x/airdrop/types"
)

// AirdropQueryClient returns the govtypes.QueryClient
// initialized with the clients grpc connection
func (c *Client) AirdropQueryClient() airdroptypes.QueryClient {
	return airdroptypes.NewQueryClient(c.grpcConn)
}

// QueryProposal sends a grpc query with the given proposalID
// and returns the govtypes.Proposal object
func (c *Client) QueryAirdropParams() (*airdroptypes.Params, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	queryResponse, err := c.AirdropQueryClient().Params(ctx, &airdroptypes.ParamsRequest{})
	if err != nil {
		return nil, err
	}
	return queryResponse.Params, nil
}
