package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AccountKeeper defines the expected interface contract defined by the x/auth module.
type AccountKeeper interface {
	NewAccount(ctx context.Context, account sdk.AccountI) sdk.AccountI
	SetAccount(ctx context.Context, account sdk.AccountI)
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	GetModuleAddress(moduleName string) sdk.AccAddress
}

// BankKeeper defines the expected interface contract defined by the x/bank module.
type BankKeeper interface {
	SendCoinsFromModuleToAccount(
		ctx context.Context,
		senderModule string,
		recipientAddr sdk.AccAddress,
		amt sdk.Coins,
	) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
}

// StakingKeeper defines the expected interface contract defined by the x/staking module.
type StakingKeeper interface {
	GetDelegatorDelegations(
		ctx context.Context,
		delegator sdk.AccAddress,
		maxRetrieve uint16,
	) ([]stakingtypes.Delegation, error)
}

type DistributionKeeper interface {
	FundCommunityPool(ctx context.Context, amount sdk.Coins, sender sdk.AccAddress) error
}
