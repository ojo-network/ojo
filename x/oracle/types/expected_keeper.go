package types

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// StakingKeeper defines the expected interface contract defined by the x/staking
// module.
type StakingKeeper interface {
	Validator(ctx context.Context, address sdk.ValAddress) (stakingtypes.ValidatorI, error)
	GetBondedValidatorsByPower(ctx context.Context) ([]stakingtypes.Validator, error)
	PowerReduction(ctx context.Context) (res sdkmath.Int)
}

// AccountKeeper defines the expected interface contract defined by the x/auth
// module.
type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx context.Context, moduleName string) sdk.ModuleAccountI

	// only used for simulation
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
}

// BankKeeper defines the expected interface contract defined by the x/bank
// module.
type BankKeeper interface {
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
}
