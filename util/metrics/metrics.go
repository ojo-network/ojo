package metrics

import "github.com/armon/go-metrics"

func RecordEndBlockMetrics(oracleKeeper OracleKeeper, ctx sdk.Context) {
	if !k.telemetryEnabled {
		return
	}

	k.IterateMissCounters(ctx, func(operator sdk.ValAddress, missCounter uint64) bool {
		metrics.SetGaugeWithLabels(
			[]string{"miss_counter"},
			float32(missCounter),
			[]metrics.Label{{Name: "address", Value: operator.String()}},
		)
		return false
	})

	medians := k.AllMedianPrices(ctx)
	medians = *medians.FilterByBlock(medians.NewestBlockNum())
	for _, median := range medians {
		metrics.SetGaugeWithLabels(
			[]string{"median"},
			float32(median.ExchangeRateTuple.ExchangeRate.MustFloat64()),
			[]metrics.Label{{Name: "denom", Value: median.ExchangeRateTuple.Denom}},
		)
	}

	medianDeviations := k.AllMedianDeviationPrices(ctx)
	medianDeviations = *medianDeviations.FilterByBlock(medianDeviations.NewestBlockNum())
	for _, medianDeviation := range medianDeviations {
		metrics.SetGaugeWithLabels(
			[]string{"median_deviation"},
			float32(medianDeviation.ExchangeRateTuple.ExchangeRate.MustFloat64()),
			[]metrics.Label{{Name: "denom", Value: medianDeviation.ExchangeRateTuple.Denom}},
		)
	}
}

func RecordExchangeRate(denom string, exchangeRate sdk.Dec) {
	metrics.SetGaugeWithLabels(
		[]string{"exchange_rate"},
		float32(exchangeRate.MustFloat64()),
		[]metrics.Label{{Name: "denom", Value: denom}},
	)
}
