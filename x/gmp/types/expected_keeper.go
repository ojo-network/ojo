package types

import (
	context "context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	gasestimatetypes "github.com/ojo-network/ojo/x/gasestimate/types"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

// OracleKeeper defines the expected Oracle interface that is needed by the gmp module.
type OracleKeeper interface {
	GetExchangeRate(ctx sdk.Context, symbol string) (math.LegacyDec, error)
	GetExponent(ctx sdk.Context, denom string) (uint32, error)
	MaximumMedianStamps(ctx sdk.Context) uint64
	HistoricMedians(ctx sdk.Context, denom string, numStamps uint64) oracletypes.PriceStamps
	HistoricDeviations(ctx sdk.Context, denom string, numStamps uint64) oracletypes.PriceStamps
}

type IBCTransferKeeper interface {
	Transfer(goCtx context.Context, msg *ibctransfertypes.MsgTransfer) (*ibctransfertypes.MsgTransferResponse, error)
}

type BankKeeper interface {
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
	GetDenomMetaData(ctx context.Context, denom string) (banktypes.Metadata, bool)
	SetDenomMetaData(ctx context.Context, denomMetaData banktypes.Metadata)
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

type GasEstimateKeeper interface {
	GetParams(ctx sdk.Context) (params gasestimatetypes.Params)
	GetGasEstimate(ctx sdk.Context, network string) (gasestimatetypes.GasEstimate, error)
}
