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

func (c *Client) QueryGovAccount() (account authtypes.ModuleAccount, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	queryResponse, err := c.AuthQueryClient().ModuleAccountByName(ctx, &authtypes.QueryModuleAccountByNameRequest{Name: "gov"})
	if err != nil {
		return
	}

	err = account.Unmarshal(queryResponse.Account.Value)
	return
}
