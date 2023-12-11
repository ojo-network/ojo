package types

import (
	context "context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

// OracleKeeper defines the expected Oracle interface that is needed by the gmp module.
type OracleKeeper interface {
	GetExchangeRate(ctx sdk.Context, symbol string) (sdk.Dec, error)
	GetExponent(ctx sdk.Context, denom string) (uint32, error)
	MaximumMedianStamps(ctx sdk.Context) uint64
	HistoricMedians(ctx sdk.Context, denom string, numStamps uint64) oracletypes.PriceStamps
	HistoricDeviations(ctx sdk.Context, denom string, numStamps uint64) oracletypes.PriceStamps
}

type IBCTransferKeeper interface {
	Transfer(goCtx context.Context, msg *ibctransfertypes.MsgTransfer) (*ibctransfertypes.MsgTransferResponse, error)
}
