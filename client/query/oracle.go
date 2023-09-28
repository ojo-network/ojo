package query

import (
	"context"

	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

// OracleQueryClient returns the oracletypes.QueryClient
// initialized with the clients grpc connection
func (c *Client) OracleQueryClient() oracletypes.QueryClient {
	return oracletypes.NewQueryClient(c.grpcConn)
}

// QueryParams returns the params from the oracle module
func (c *Client) QueryParams() (oracletypes.Params, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	queryResponse, err := c.OracleQueryClient().Params(ctx, &oracletypes.QueryParams{})
	if err != nil {
		return oracletypes.Params{}, err
	}
	return queryResponse.Params, nil
}

// QueryExchangeRates returns the exchange rates from the oracle module
func (c *Client) QueryExchangeRates() ([]oracletypes.PriceStamp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	queryResponse, err := c.OracleQueryClient().ExchangeRates(ctx, &oracletypes.QueryExchangeRates{})
	if err != nil {
		return nil, err
	}
	return queryResponse.ExchangeRates, nil
}

// QueryMedians returns the medians from the oracle module
func (c *Client) QueryMedians() (oracletypes.PriceStamps, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	queryResponse, err := c.OracleQueryClient().Medians(ctx, &oracletypes.QueryMedians{})
	if err != nil {
		return nil, err
	}
	return queryResponse.Medians, nil
}

func (c *Client) QueryMedianDeviations() (oracletypes.PriceStamps, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	queryResponse, err := c.OracleQueryClient().MedianDeviations(ctx, &oracletypes.QueryMedianDeviations{})
	if err != nil {
		return nil, err
	}
	return queryResponse.MedianDeviations, nil
}
