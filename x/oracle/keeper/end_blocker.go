package keeper

import (
	"sort"

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
	allAssetPrice := k.GetAllAssetPrices(ctx, asset)
	total := len(allAssetPrice)

	sort.Slice(allAssetPrice, func(i, j int) bool {
		return allAssetPrice[i].Timestamp < allAssetPrice[j].Timestamp
	})

	for i, price := range allAssetPrice {
		// We don't remove the last element
		if i < total-1 {
			k.RemovePrice(ctx, price.Asset, price.Timestamp)
		}
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
