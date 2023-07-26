package client

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/ojo-network/ojo/client/query"
	"github.com/ojo-network/ojo/client/tx"
	"github.com/rs/zerolog"
)

// OjoClient is a helper for initializing a keychain, a cosmos-sdk client context,
// and sending transactions/queries to a specific Ojo node
type OjoClient struct {
	QueryClient *query.Client
	TxClient    *tx.Client
}

// NewOjoClient returns a new instance of the OjoClient with initialized
// query and transaction clients
func NewOjoClient(
	chainID string,
	tmrpcEndpoint string,
	grpcEndpoint string,
	accountName string,
	accountMnemonic string,
) (oc *OjoClient, err error) {
	oc = &OjoClient{}
	oc.QueryClient, err = query.NewClient(grpcEndpoint)
	if err != nil {
		return nil, err
	}
	oc.TxClient, err = tx.NewClient(chainID, tmrpcEndpoint, accountName, accountMnemonic)
	return oc, err
}

// NewChainHeight returns a new instance of the ChainHeight struct
// using the OjoClient's transaction sdk.client
func (oc *OjoClient) NewChainHeight(ctx context.Context, logger zerolog.Logger) (*ChainHeight, error) {
	return NewChainHeight(ctx, oc.TxClient.ClientContext.Client, logger)
}

func (oc *OjoClient) QueryTxHash(hash string) (*sdk.TxResponse, error) {
	return authtx.QueryTx(*oc.TxClient.ClientContext, hash)
}
