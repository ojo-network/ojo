package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/x/oracle/types"
)

// ScheduleParamUpdatePlan schedules a param update plan.
func (k Keeper) ScheduleParamUpdatePlan(ctx sdk.Context, plan types.ParamUpdatePlan) error {
	if plan.Height < ctx.BlockHeight() {
		return types.ErrInvalidRequest.Wrap("param update cannot be scheduled in the past")
	}
	if err := k.ValidateParamChanges(ctx, plan.Keys, plan.Changes); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&plan)
	store.Set(types.KeyParamUpdatePlan(), bz)

	return nil
}

// ClearParamUpdatePlan will clear an upcoming param update plan if one exists and return
// an error if one isn't found.
func (k Keeper) ClearParamUpdatePlan(ctx sdk.Context) error {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyParamUpdatePlan())
	if bz == nil {
		return types.ErrInvalidRequest.Wrap("No param update plan found")
	}

	store.Delete(types.KeyParamUpdatePlan())
	return nil
}

// GetParamUpdatePlan will return whether an upcoming param update plan exists and the plan
// if it does.
func (k Keeper) GetParamUpdatePlan(ctx sdk.Context) (plan types.ParamUpdatePlan, havePlan bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyParamUpdatePlan())
	if bz == nil {
		return plan, false
	}

	k.cdc.MustUnmarshal(bz, &plan)
	return plan, true
}

// ValidateParamChanges validates parameter changes against the existing oracle parameters.
func (k Keeper) ValidateParamChanges(ctx sdk.Context, keys []string, changes types.Params) error {
	params := k.GetParams(ctx)

	for _, key := range keys {
		switch key {
		case string(types.KeyVotePeriod):
			params.VotePeriod = changes.VotePeriod

		case string(types.KeyVoteThreshold):
			params.VoteThreshold = changes.VoteThreshold

		case string(types.KeyRewardBands):
			params.RewardBands = changes.RewardBands

		case string(types.KeyRewardDistributionWindow):
			params.RewardDistributionWindow = changes.RewardDistributionWindow

		case string(types.KeyAcceptList):
			params.AcceptList = changes.AcceptList.Normalize()

		case string(types.KeyMandatoryList):
			params.MandatoryList = changes.MandatoryList.Normalize()

		case string(types.KeySlashFraction):
			params.SlashFraction = changes.SlashFraction

		case string(types.KeySlashWindow):
			params.SlashWindow = changes.SlashWindow

		case string(types.KeyMinValidPerWindow):
			params.MinValidPerWindow = changes.MinValidPerWindow

		case string(types.KeyHistoricStampPeriod):
			params.HistoricStampPeriod = changes.HistoricStampPeriod

		case string(types.KeyMedianStampPeriod):
			params.MedianStampPeriod = changes.MedianStampPeriod

		case string(types.KeyMaximumPriceStamps):
			params.MaximumPriceStamps = changes.MaximumPriceStamps

		case string(types.KeyMaximumMedianStamps):
			params.MaximumMedianStamps = changes.MaximumMedianStamps
		}
	}

	return params.Validate()
}

// ExecuteParamUpdatePlan will execute a given param update plan.
func (k Keeper) ExecuteParamUpdatePlan(ctx sdk.Context, plan types.ParamUpdatePlan) {
	for _, key := range plan.Keys {
		switch key {
		case string(types.KeyVotePeriod):
			k.SetVotePeriod(ctx, plan.Changes.VotePeriod)

		case string(types.KeyVoteThreshold):
			k.SetVoteThreshold(ctx, plan.Changes.VoteThreshold)

		case string(types.KeyRewardBands):
			k.SetRewardBand(ctx, plan.Changes.RewardBands)

		case string(types.KeyRewardDistributionWindow):
			k.SetRewardDistributionWindow(ctx, plan.Changes.RewardDistributionWindow)

		case string(types.KeyAcceptList):
			k.SetAcceptList(ctx, plan.Changes.AcceptList.Normalize())

		case string(types.KeyMandatoryList):
			k.SetMandatoryList(ctx, plan.Changes.MandatoryList.Normalize())

		case string(types.KeySlashFraction):
			k.SetSlashFraction(ctx, plan.Changes.SlashFraction)

		case string(types.KeySlashWindow):
			k.SetSlashWindow(ctx, plan.Changes.SlashWindow)

		case string(types.KeyMinValidPerWindow):
			k.SetMinValidPerWindow(ctx, plan.Changes.MinValidPerWindow)

		case string(types.KeyHistoricStampPeriod):
			k.SetHistoricStampPeriod(ctx, plan.Changes.HistoricStampPeriod)

		case string(types.KeyMedianStampPeriod):
			k.SetMedianStampPeriod(ctx, plan.Changes.MedianStampPeriod)

		case string(types.KeyMaximumPriceStamps):
			k.SetMaximumPriceStamps(ctx, plan.Changes.MaximumPriceStamps)

		case string(types.KeyMaximumMedianStamps):
			k.SetMaximumMedianStamps(ctx, plan.Changes.MaximumMedianStamps)
		}
	}
}
