package airdrop

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/util"
	"github.com/ojo-network/ojo/x/airdrop/keeper"
	"github.com/ojo-network/ojo/x/airdrop/types"
)

const (
	BatchSize = 10
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx context.Context, k keeper.Keeper) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	createOriginAccounts(sdkCtx, k)
	return distributeUnclaimedAirdrops(sdkCtx, k)
}

// createOriginAccounts creates the airdrop accounts for all the
// addresses in the airdrop module account.
func createOriginAccounts(ctx sdk.Context, k keeper.Keeper) {
	airdropAccounts := k.PaginatedAirdropAccounts(ctx, types.AirdropAccount_STATE_CREATED, BatchSize)
	for _, airdropAccount := range airdropAccounts {
		err := k.CreateAirdropAccount(ctx, airdropAccount)
		if err != nil {
			ctx.Logger().Error("error creating airdrop account", err)
		}
	}
}

func distributeUnclaimedAirdrops(ctx sdk.Context, k keeper.Keeper) error {
	if ctx.BlockHeight() < util.SafeUint64ToInt64(k.GetParams(ctx).ExpiryBlock) {
		return nil
	}

	for _, aa := range k.PaginatedAirdropAccounts(ctx, types.AirdropAccount_STATE_UNCLAIMED, BatchSize) {
		if aa.VerifyNotClaimed() != nil {
			continue
		}
		err := distributeUnclaimedAirdrop(ctx, k, aa)
		if err != nil {
			return err
		}
	}
	return nil
}

func distributeUnclaimedAirdrop(ctx sdk.Context, k keeper.Keeper, aa *types.AirdropAccount) error {
	k.SetClaimAmount(ctx, aa)
	err := k.MintClaimTokensToDistribution(ctx, aa)
	if err != nil {
		return err
	}
	aa.ClaimAddress = k.DistributionModuleAddress(ctx).String()
	return k.ChangeAirdropAccountState(ctx, aa, types.AirdropAccount_STATE_CLAIMED)
}
