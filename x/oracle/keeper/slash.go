package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/util"
	"github.com/ojo-network/ojo/x/oracle/types"
)

// SlashAndResetMissCounters iterates over all the current missed counters and
// calculates the "valid vote rate" as:
// (possibleWinsPerSlashWindow - missCounter)/possibleWinsPerSlashWindow.
//
// If the valid vote rate is below the minValidPerWindow, the validator will be
// slashed and jailed.
func (k Keeper) SlashAndResetMissCounters(ctx sdk.Context) {
	var (
		possibleWinsPerSlashWindow = k.PossibleWinsPerSlashWindow(ctx)
		minValidPerWindow          = k.MinValidPerWindow(ctx)

		distributionHeight = ctx.BlockHeight() - sdk.ValidatorUpdateDelay - 1
		slashFraction      = k.SlashFraction(ctx)
		powerReduction     = k.StakingKeeper.PowerReduction(ctx)
	)

	k.IterateMissCounters(ctx, func(operator sdk.ValAddress, missCounter uint64) bool {
		validVotes := sdk.NewInt(possibleWinsPerSlashWindow - int64(missCounter))
		validVoteRate := sdk.NewDecFromInt(validVotes).QuoInt64(possibleWinsPerSlashWindow)

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

// PossibleWinsPerSlashWindow returns the total number of possible correct votes
// that a validator can have per asset multiplied by the number of vote
// periods in the slash window
func (k Keeper) PossibleWinsPerSlashWindow(ctx sdk.Context) int64 {
	slashWindow := int64(k.SlashWindow(ctx))
	votePeriod := int64(k.VotePeriod(ctx))

	votePeriodsPerWindow := sdk.NewDec(slashWindow).QuoInt64(votePeriod).TruncateInt64()
	numberOfAssets := int64(len(k.GetParams(ctx).AcceptList))

	return (votePeriodsPerWindow * numberOfAssets)
}

// SetValidatorRewardSet will take all the current validators and store them
// in the ValidatorRewardSet to earn rewwards in the current Slash Window.
func (k Keeper) SetValidatorRewardSet(ctx sdk.Context) {
	validatorRewardSet := types.ValidatorRewardSet{
		ValidatorMap: make(map[string]bool),
	}
	for _, v := range k.StakingKeeper.GetBondedValidatorsByPower(ctx) {
		addr := v.GetOperator()
		validatorRewardSet.ValidatorMap[addr.String()] = true
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&validatorRewardSet)
	store.Set(types.KeyValidatorRewardSet(uint64(ctx.BlockHeight())), bz)
}

// CurrentValidatorRewardSet returns the latest ValidatorRewardSet in the store.
func (k Keeper) CurrentValidatorRewardSet(
	ctx sdk.Context,
	handler func(types.ValidatorRewardSet) bool,
) {
	store := ctx.KVStore(k.storeKey)

	// make sure we have one zero byte to correctly seperate blocknum
	prefix := util.ConcatBytes(1, types.KeyPrefixValidatorRewardSet)
	iter := sdk.KVStoreReversePrefixIteratorPaginated(store, prefix, 1, 1)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		validatorRewardSet := types.ValidatorRewardSet{
			ValidatorMap: make(map[string]bool),
		}
		k.cdc.MustUnmarshal(iter.Value(), &validatorRewardSet)
		if handler(validatorRewardSet) {
			break
		}
	}
}
