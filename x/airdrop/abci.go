package airdrop

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/x/airdrop/keeper"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) error {
	if ctx.BlockHeight() == int64(k.GetParams(ctx).ExpiryBlock) {
		for _, aa := range k.GetAllAirdropAccounts(ctx) {
			if aa.VerifyNotClaimed() == nil {
				k.SetClaimAmount(ctx, aa)
				k.MintClaimTokensToAirdrop(ctx, aa)
				aa.ClaimAddress = k.DistributionModuleAddress(ctx).String()
				k.SetAirdropAccount(ctx, aa)
			}
		}
	}
	return nil
}
