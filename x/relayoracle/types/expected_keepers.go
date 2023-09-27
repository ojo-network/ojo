package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"

	oracleTypes "github.com/ojo-network/ojo/x/oracle/types"
)

type OracleKeeper interface {
	GetExchangeRate(ctx sdk.Context, denom string) (sdk.Dec, error)
	IterateExchangeRatesWithDenoms(ctx sdk.Context, denoms []string, blocknum uint64) (oracleTypes.PriceStamps, error)
	IterateHistoricPricesForDenoms(ctx sdk.Context, prefix []byte, denoms []string, numStamps uint) oracleTypes.PriceStamps
	HasActiveExchangeRate(ctx sdk.Context, denom []string) bool
	MaximumMedianStamps(ctx sdk.Context) (res uint64)
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
