package query

import (
	"context"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AuthQueryClient returns the authtypes.QueryClient
// initialized with the clients grpc connection
func (c *Client) AuthQueryClient() authtypes.QueryClient {
	return authtypes.NewQueryClient(c.grpcConn)
}

// QueryGovAccount returns the gov module account for the chain
func (c *Client) QueryGovAccount() (account authtypes.ModuleAccount, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	request := &authtypes.QueryModuleAccountByNameRequest{Name: "gov"}
	queryResponse, err := c.AuthQueryClient().ModuleAccountByName(ctx, request)
	if err != nil {
		return
	}

	err = account.Unmarshal(queryResponse.Account.Value)
	return
}
