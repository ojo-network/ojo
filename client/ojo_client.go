package client

import (
	"context"

	"github.com/ojo-network/ojo/client/query"
	"github.com/ojo-network/ojo/client/tx"
	"github.com/rs/zerolog"
)

// OjoClient is a helper for initializing a keychain, a cosmos-sdk client context,
// and sending transactions/queries to a specific Umee node
type OjoClient struct {
	QueryClient *query.Client
	TxClient    *tx.Client
}

func NewOjoClient(
	chainID string,
	tmrpcEndpoint string,
	grpcEndpoint string,
	accountName string,
	accountMnemonic string,
) (uc *OjoClient, err error) {
	uc = &OjoClient{}
	uc.QueryClient, err = query.NewQueryClient(grpcEndpoint)
	if err != nil {
		return nil, err
	}
	uc.TxClient, err = tx.NewTxClient(chainID, tmrpcEndpoint, accountName, accountMnemonic)
	return uc, err
}

func (oc *OjoClient) NewChainHeight(ctx context.Context, logger zerolog.Logger) (*ChainHeight, error) {
	return NewChainHeight(ctx, oc.TxClient.ClientContext.Client, logger)
}
