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
	GetValidator(ctx context.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, err error)
	GetValidatorByConsAddr(ctx context.Context, consAddr sdk.ConsAddress) (validator stakingtypes.Validator, err error)
	AddValidatorTokensAndShares(ctx context.Context, validator stakingtypes.Validator,
		tokensToAdd sdkmath.Int,
	) (valOut stakingtypes.Validator, addedShares sdkmath.LegacyDec, err error)
	RemoveValidatorTokens(ctx context.Context,
		validator stakingtypes.Validator, tokensToRemove sdkmath.Int,
	) (stakingtypes.Validator, error)
	PowerReduction(ctx context.Context) (res sdkmath.Int)
}
