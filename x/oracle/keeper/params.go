package keeper

import (
	"fmt"
	"strings"

	"github.com/ojo-network/ojo/x/oracle/types"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// VotePeriod returns the number of blocks during which voting takes place.
func (k Keeper) VotePeriod(ctx sdk.Context) uint64 {
	params := k.GetParams(ctx)
	return params.VotePeriod
}

// SetVotePeriod updates the number of blocks during which voting takes place.
func (k Keeper) SetVotePeriod(ctx sdk.Context, votePeriod uint64) {
	params := k.GetParams(ctx)
	params.VotePeriod = votePeriod
	k.SetParams(ctx, params)
}

// VoteThreshold returns the minimum percentage of votes that must be received
// for a ballot to pass.
func (k Keeper) VoteThreshold(ctx sdk.Context) math.LegacyDec {
	params := k.GetParams(ctx)
	return params.VoteThreshold
}

// SetVoteThreshold updates the minimum percentage of votes that must be received
// for a ballot to pass.
func (k Keeper) SetVoteThreshold(ctx sdk.Context, voteThreshold math.LegacyDec) {
	params := k.GetParams(ctx)
	params.VoteThreshold = voteThreshold
	k.SetParams(ctx, params)
}

// RewardBand returns the ratio of allowable exchange rate error that a validator
// can be rewarded.
func (k Keeper) RewardBands(ctx sdk.Context) types.RewardBandList {
	params := k.GetParams(ctx)
	return params.RewardBands
}

// VoteThreshold updates the ratio of allowable exchange rate error that a validator
// can be rewarded.
func (k Keeper) SetRewardBand(ctx sdk.Context, rewardBands types.RewardBandList) {
	params := k.GetParams(ctx)
	params.RewardBands = rewardBands
	k.SetParams(ctx, params)
}

// RewardDistributionWindow returns the number of vote periods during which
// seigniorage reward comes in and then is distributed.
func (k Keeper) RewardDistributionWindow(ctx sdk.Context) uint64 {
	params := k.GetParams(ctx)
	return params.RewardDistributionWindow
}

// SetRewardDistributionWindow updates the number of vote periods during which
// seigniorage reward comes in and then is distributed.
func (k Keeper) SetRewardDistributionWindow(ctx sdk.Context, rewardDistributionWindow uint64) {
	params := k.GetParams(ctx)
	params.RewardDistributionWindow = rewardDistributionWindow
	k.SetParams(ctx, params)
}

// AcceptList returns the denom list that can be activated
func (k Keeper) AcceptList(ctx sdk.Context) types.DenomList {
	params := k.GetParams(ctx)
	return params.AcceptList
}

// SetAcceptList updates the accepted list of assets supported by the x/oracle
// module.
func (k Keeper) SetAcceptList(ctx sdk.Context, acceptList types.DenomList) {
	params := k.GetParams(ctx)
	params.AcceptList = acceptList
	k.SetParams(ctx, params)
}

// MandatoryList returns the denom list that are mandatory
func (k Keeper) MandatoryList(ctx sdk.Context) types.DenomList {
	params := k.GetParams(ctx)
	return params.MandatoryList
}

// SetMandatoryList updates the mandatory list of assets supported by the x/oracle
// module.
func (k Keeper) SetMandatoryList(ctx sdk.Context, mandatoryList types.DenomList) {
	params := k.GetParams(ctx)
	params.MandatoryList = mandatoryList
	k.SetParams(ctx, params)
}

// SlashFraction returns the oracle voting penalty rate.
func (k Keeper) SlashFraction(ctx sdk.Context) math.LegacyDec {
	params := k.GetParams(ctx)
	return params.SlashFraction
}

// SetSlashFraction updates the oracle voting penalty rate.
func (k Keeper) SetSlashFraction(ctx sdk.Context, slashFraction math.LegacyDec) {
	params := k.GetParams(ctx)
	params.SlashFraction = slashFraction
	k.SetParams(ctx, params)
}

// SlashWindow returns the number of total blocks in a slash window.
func (k Keeper) SlashWindow(ctx sdk.Context) uint64 {
	params := k.GetParams(ctx)
	return params.SlashWindow
}

// SetSlashWindow updates the number of total blocks in a slash window.
func (k Keeper) SetSlashWindow(ctx sdk.Context, slashWindow uint64) {
	params := k.GetParams(ctx)
	params.SlashWindow = slashWindow
	k.SetParams(ctx, params)
}

// MinValidPerWindow returns the oracle slashing threshold.
func (k Keeper) MinValidPerWindow(ctx sdk.Context) math.LegacyDec {
	params := k.GetParams(ctx)
	return params.MinValidPerWindow
}

// MinValidPerWindow updates the oracle slashing threshold.
func (k Keeper) SetMinValidPerWindow(ctx sdk.Context, minValidPerWindow math.LegacyDec) {
	params := k.GetParams(ctx)
	params.MinValidPerWindow = minValidPerWindow
	k.SetParams(ctx, params)
}

// GetParams returns the total set of oracle parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.KeyParams)
	if b == nil {
		return
	}

	k.cdc.MustUnmarshal(b, &params)
	return
}

// SetParams sets the total set of oracle parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := k.cdc.MustMarshal(&params)
	store.Set(types.KeyParams, b)
}

// HistoricStampPeriod returns the amount of blocks the oracle module waits
// before recording a new historic price.
func (k Keeper) HistoricStampPeriod(ctx sdk.Context) uint64 {
	params := k.GetParams(ctx)
	return params.HistoricStampPeriod
}

// SetHistoricStampPeriod updates the amount of blocks the oracle module waits
// before recording a new historic price.
func (k Keeper) SetHistoricStampPeriod(ctx sdk.Context, historicPriceStampPeriod uint64) {
	params := k.GetParams(ctx)
	params.HistoricStampPeriod = historicPriceStampPeriod
	k.SetParams(ctx, params)
}

// MedianStampPeriod returns the amount blocks the oracle module waits between
// calculating a new median and standard deviation of that median.
func (k Keeper) MedianStampPeriod(ctx sdk.Context) uint64 {
	params := k.GetParams(ctx)
	return params.MedianStampPeriod
}

// SetMedianStampPeriod updates the amount blocks the oracle module waits between
// calculating a new median and standard deviation of that median.
func (k Keeper) SetMedianStampPeriod(ctx sdk.Context, medianStampPeriod uint64) {
	params := k.GetParams(ctx)
	params.MedianStampPeriod = medianStampPeriod
	k.SetParams(ctx, params)
}

// MaximumPriceStamps returns the maximum amount of historic prices the oracle
// module will hold.
func (k Keeper) MaximumPriceStamps(ctx sdk.Context) uint64 {
	params := k.GetParams(ctx)
	return params.MaximumPriceStamps
}

// SetMaximumPriceStamps updates the the maximum amount of historic prices the
// oracle module will hold.
func (k Keeper) SetMaximumPriceStamps(ctx sdk.Context, maximumPriceStamps uint64) {
	params := k.GetParams(ctx)
	params.MaximumPriceStamps = maximumPriceStamps
	k.SetParams(ctx, params)
}

// MaximumMedianStamps returns the maximum amount of medians the oracle module will
// hold.
func (k Keeper) MaximumMedianStamps(ctx sdk.Context) uint64 {
	params := k.GetParams(ctx)
	return params.MaximumMedianStamps
}

// SetMaximumMedianStamps updates the the maximum amount of medians the oracle module will
// hold.
func (k Keeper) SetMaximumMedianStamps(ctx sdk.Context, maximumMedianStamps uint64) {
	params := k.GetParams(ctx)
	params.MaximumMedianStamps = maximumMedianStamps
	k.SetParams(ctx, params)
}

// PriceExpiryTime returns the expiry in unix time for elys prices.
func (k Keeper) PriceExpiryTime(ctx sdk.Context) uint64 {
	params := k.GetParams(ctx)
	return params.PriceExpiryTime
}

// SetPriceExpiryTime updates the expiry in unix time for elys prices.
func (k Keeper) SetPriceExpiryTime(ctx sdk.Context, priceExpiryTime uint64) {
	params := k.GetParams(ctx)
	params.PriceExpiryTime = priceExpiryTime
	k.SetParams(ctx, params)
}

// LifeTimeInBlocks returns the life time of an elys price in blocks.
func (k Keeper) LifeTimeInBlocks(ctx sdk.Context) uint64 {
	params := k.GetParams(ctx)
	return params.LifeTimeInBlocks
}

// SetLifeTimeInBlocks updates the life time of an elys price in blocks.
func (k Keeper) SetLifeTimeInBlocks(ctx sdk.Context, lifeTimeInBlocks uint64) {
	params := k.GetParams(ctx)
	params.LifeTimeInBlocks = lifeTimeInBlocks
	k.SetParams(ctx, params)
}

// CurrencyPairProviders returns the current Currency Pair Providers the price feeder
// will query when starting up.
func (k Keeper) CurrencyPairProviders(ctx sdk.Context) types.CurrencyPairProvidersList {
	params := k.GetParams(ctx)
	return params.CurrencyPairProviders
}

// SetCurrencyPairProviders updates the current Currency Pair Providers the price feeder
// will query when starting up.
func (k Keeper) SetCurrencyPairProviders(
	ctx sdk.Context,
	currencyPairProviders types.CurrencyPairProvidersList,
) {
	params := k.GetParams(ctx)
	params.CurrencyPairProviders = currencyPairProviders
	k.SetParams(ctx, params)
}

// CurrencyDeviationThresholds returns the current Currency Deviation Thesholds the
// price feeder will query when starting up.
func (k Keeper) CurrencyDeviationThresholds(ctx sdk.Context) types.CurrencyDeviationThresholdList {
	params := k.GetParams(ctx)
	return params.CurrencyDeviationThresholds
}

// SetCurrencyDeviationThresholds updates the current Currency Deviation Thesholds the
// price feeder will query when starting up.
func (k Keeper) SetCurrencyDeviationThresholds(
	ctx sdk.Context,
	currencyDeviationThresholds types.CurrencyDeviationThresholdList,
) {
	params := k.GetParams(ctx)
	params.CurrencyDeviationThresholds = currencyDeviationThresholds
	k.SetParams(ctx, params)
}

func (k Keeper) GetExponent(ctx sdk.Context, denom string) (uint32, error) {
	params := k.GetParams(ctx)
	for _, v := range params.AcceptList {
		if strings.EqualFold(v.SymbolDenom, denom) {
			return v.Exponent, nil
		}
	}
	return 0, fmt.Errorf("unable to find exponent for %s", denom)
}
