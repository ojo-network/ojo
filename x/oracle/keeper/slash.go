package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SlashAndResetMissCounters iterates over all the current missed counters and
// calculates the "valid vote rate" as:
// (votePeriodsPerWindow - missCounter)/votePeriodsPerWindow.
//
// If the valid vote rate is below the minValidPerWindow, the validator will be
// slashed and jailed.
func (k Keeper) SlashAndResetMissCounters(ctx sdk.Context) {
	var (
		height             = ctx.BlockHeight()
		distributionHeight = height - sdk.ValidatorUpdateDelay - 1

		slashWindow          = int64(k.SlashWindow(ctx))
		votePeriod           = int64(k.VotePeriod(ctx))
		votePeriodsPerWindow = sdk.NewDec(slashWindow).QuoInt64(votePeriod).TruncateInt64()

		numberOfAssets              = int64(len(k.GetParams(ctx).AcceptList))
		possibleVotesPerSlashWindow = votePeriodsPerWindow * numberOfAssets
		minValidPerWindow           = k.MinValidPerWindow(ctx)

		slashFraction  = k.SlashFraction(ctx)
		powerReduction = k.StakingKeeper.PowerReduction(ctx)
	)

	k.IterateMissCounters(ctx, func(operator sdk.ValAddress, missCounter uint64) bool {
		numValidVotes := sdk.NewInt(possibleVotesPerSlashWindow - int64(missCounter))
		validVoteRate := sdk.NewDecFromInt(numValidVotes).QuoInt64(possibleVotesPerSlashWindow)

		// Slash and jail the validator if their valid vote rate is smaller than the
		// minimum threshold.
		if validVoteRate.LT(minValidPerWindow) {
			validator := k.StakingKeeper.Validator(ctx, operator)
			if validator.IsBonded() && !validator.IsJailed() {
				consAddr, err := validator.GetConsAddr()
				if err != nil {
					panic(err)
				}

				k.StakingKeeper.Slash(
					ctx,
					consAddr,
					distributionHeight,
					validator.GetConsensusPower(powerReduction), slashFraction,
				)

				k.StakingKeeper.Jail(ctx, consAddr)
			}
		}

		k.DeleteMissCounter(ctx, operator)
		return false
	})
}
