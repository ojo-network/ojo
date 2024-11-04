package query

import (
	"context"

	gmptypes "github.com/ojo-network/ojo/x/gmp/types"
)

func (c *Client) GmpQueryClient() gmptypes.QueryClient {
	return gmptypes.NewQueryClient(c.grpcConn)
}

func (c *Client) QueryPayments() (*gmptypes.AllPaymentsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	queryResponse, err := c.GmpQueryClient().AllPayments(
		ctx,
		&gmptypes.AllPaymentsRequest{},
	)

	return queryResponse, err
}
