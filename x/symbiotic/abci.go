package symbiotic

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/util"
	"github.com/ojo-network/ojo/x/symbiotic/keeper"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx context.Context, k keeper.Keeper) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(sdkCtx)
	blockHeight := util.SafeInt64ToUint64(sdkCtx.BlockHeight())
	syncPeriod := util.SafeInt64ToUint64(params.SymbioticSyncPeriod)

	if util.IsPeriodLastBlock(sdkCtx, syncPeriod) {
		prunePeriod := params.MaximumCachedBlockHashes * syncPeriod
		if prunePeriod < blockHeight {
			k.PruneBlockHashesBeforeBlock(sdkCtx, blockHeight-prunePeriod)
		}
	}

	if err := k.SymbioticUpdateValidatorsPower(ctx); err != nil {
		k.Logger(sdkCtx).With(err).Error("Symbiotic val update error")
	}

	return nil
}
