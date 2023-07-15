package query

import (
	"context"

	"cosmossdk.io/math"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// GovQueryClient returns the govtypes.QueryClient
// initialized with the clients grpc connection
func (c *Client) BankQueryClient() banktypes.QueryClient {
	return banktypes.NewQueryClient(c.grpcConn)
}

func (c *Client) QueryBalance(address string, denom string) (math.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	queryResponse, err := c.BankQueryClient().Balance(
		ctx,
		&banktypes.QueryBalanceRequest{Address: address, Denom: denom},
	)
	if err != nil {
		return math.Int{}, err
	}
	return queryResponse.Balance.Amount, nil
}
