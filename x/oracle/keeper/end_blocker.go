package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ojo-network/ojo/util"
	"github.com/ojo-network/ojo/util/metrics"
)

// PruneAllPrices deletes all historic prices, medians, and median deviations
// outside pruning period determined by the stamp period multiplied by the maximum stamps.
func (k *Keeper) PruneAllPrices(ctx sdk.Context) {
	params := k.GetParams(ctx)
	blockHeight := util.SafeInt64ToUint64(ctx.BlockHeight())

	if k.IsPeriodLastBlock(ctx, params.HistoricStampPeriod) {
		pruneHistoricPeriod := params.HistoricStampPeriod * params.MaximumPriceStamps
		if pruneHistoricPeriod < blockHeight {
			k.PruneHistoricPricesBeforeBlock(ctx, blockHeight-pruneHistoricPeriod)
		}

		if k.IsPeriodLastBlock(ctx, params.MedianStampPeriod) {
			pruneMedianPeriod := params.MedianStampPeriod * params.MaximumMedianStamps
			if pruneMedianPeriod < blockHeight {
				k.PruneMediansBeforeBlock(ctx, blockHeight-pruneMedianPeriod)
				k.PruneMedianDeviationsBeforeBlock(ctx, blockHeight-pruneMedianPeriod)
			}
		}
	}
}

// PruneElysPrices prunes elys prices for a given asset except the latest one.
func (k *Keeper) PruneElysPrices(ctx sdk.Context, asset string) {
	prices := k.GetAllAssetPrices(ctx, asset)
	if len(prices) <= 1 {
		return // nothing to prune
	}

	// Find the newest price
	latestIdx, latestTs := 0, prices[0].Timestamp
	for i := 1; i < len(prices); i++ {
		if prices[i].Timestamp > latestTs {
			latestIdx, latestTs = i, prices[i].Timestamp
		}
	}

	// Remove everything except the newest
	for i, p := range prices {
		if i == latestIdx {
			continue
		}
		k.RemovePrice(ctx, p.Asset, p.Timestamp)
	}
}

// IsPeriodLastBlock returns true if we are at the last block of the period
func (k *Keeper) IsPeriodLastBlock(ctx sdk.Context, blocksPerPeriod uint64) bool {
	return (util.SafeInt64ToUint64(ctx.BlockHeight())+1)%blocksPerPeriod == 0
}

// RecordEndBlockMetrics records miss counter and price metrics at the end of the block
func (k *Keeper) RecordEndBlockMetrics(ctx sdk.Context) {
	if !k.telemetryEnabled {
		return
	}

	k.IterateMissCounters(ctx, func(operator string, missCounter uint64) bool {
		metrics.RecordMissCounter(operator, missCounter)
		return false
	})

	medians := k.AllMedianPrices(ctx)
	medians = *medians.NewestPrices()
	for _, median := range medians {
		metrics.RecordMedianPrice(median.ExchangeRate.Denom, median.ExchangeRate.Amount)
	}

	medianDeviations := k.AllMedianDeviationPrices(ctx)
	medianDeviations = *medianDeviations.NewestPrices()
	for _, medianDeviation := range medianDeviations {
		metrics.RecordMedianDeviationPrice(medianDeviation.ExchangeRate.Denom, medianDeviation.ExchangeRate.Amount)
	}
}
