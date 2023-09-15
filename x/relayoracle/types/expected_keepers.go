package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"

	oracleTypes "github.com/ojo-network/ojo/x/oracle/types"
)

type OracleKeeper interface {
	GetExchangeRate(ctx sdk.Context, denom string) (sdk.Dec, error)
	HistoricMedianDeviation(ctx sdk.Context, denom string) (*oracleTypes.PriceStamp, error)
	HistoricMedians(ctx sdk.Context, denom string, numStamps uint64) oracleTypes.PriceStamps
	HasActiveExchangeRate(ctx sdk.Context, denom string) bool
}

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}
