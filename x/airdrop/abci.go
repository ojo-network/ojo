package airdrop

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/x/airdrop/keeper"
	"github.com/ojo-network/ojo/x/airdrop/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) error {
	if ctx.BlockHeight() == int64(k.GetParams(ctx).ExpiryBlock) {
		for _, aa := range k.GetAllAirdropAccounts(ctx) {
			if aa.VerifyNotClaimed() != nil {
				continue
			}
			err := DistributeUnclaimedAirdrop(ctx, k, aa)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func DistributeUnclaimedAirdrop(ctx sdk.Context, k keeper.Keeper, aa *types.AirdropAccount) error {
	k.SetClaimAmount(ctx, aa)
	err := k.MintClaimTokensToDistribution(ctx, aa)
	if err != nil {
		return err
	}
	aa.ClaimAddress = k.DistributionModuleAddress(ctx).String()
	err = k.SetAirdropAccount(ctx, aa)
	if err != nil {
		return err
	}
	return nil
}
