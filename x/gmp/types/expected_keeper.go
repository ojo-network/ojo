package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// OracleKeeper defines the expected Oracle interface that is needed by the gmp module.
type OracleKeeper interface {
	GetExchangeRate(ctx sdk.Context, symbol string) (sdk.Dec, error)
	GetExponent(ctx sdk.Context, denom string) (uint32, error)
}
