package query

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

// OracleQueryClient returns the oracletypes.QueryClient
// initialized with the clients grpc connection
func (qc *QueryClient) OracleQueryClient() oracletypes.QueryClient {
	return oracletypes.NewQueryClient(qc.grpcConn)
}

// QueryParams returns the params from the oracle module
func (qc *QueryClient) QueryParams() (oracletypes.Params, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	queryResponse, err := qc.OracleQueryClient().Params(ctx, &oracletypes.QueryParams{})
	if err != nil {
		return oracletypes.Params{}, err
	}
	return queryResponse.Params, nil
}

// QueryExchangeRates returns the exchange rates from the oracle module
func (qc *QueryClient) QueryExchangeRates() ([]sdk.DecCoin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	queryResponse, err := qc.OracleQueryClient().ExchangeRates(ctx, &oracletypes.QueryExchangeRates{})
	if err != nil {
		return nil, err
	}
	return queryResponse.ExchangeRates, nil
}

// QueryMedians returns the medians from the oracle module
func (qc *QueryClient) QueryMedians() (sdk.DecCoins, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	queryResponse, err := qc.OracleQueryClient().Medians(ctx, &oracletypes.QueryMedians{})
	if err != nil {
		return nil, err
	}
	return queryResponse.Medians, nil
}

func (qc *QueryClient) QueryMedianDeviations() (sdk.DecCoins, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	queryResponse, err := qc.OracleQueryClient().MedianDeviations(ctx, &oracletypes.QueryMedianDeviations{})
	if err != nil {
		return nil, err
	}
	return queryResponse.MedianDeviations, nil
}
