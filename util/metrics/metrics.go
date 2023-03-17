package metrics

import (
	"github.com/armon/go-metrics"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RecordMissCounter records the miss counter gauge for a validator
func RecordMissCounter(operator sdk.ValAddress, missCounter uint64) {
	metrics.SetGaugeWithLabels(
		[]string{"miss_counter"},
		float32(missCounter),
		[]metrics.Label{{Name: "address", Value: operator.String()}},
	)
}

// RecordExchangeRate records the exchange rate gauge for a denom
func RecordExchangeRate(denom string, exchangeRate sdk.Dec) {
	metrics.SetGaugeWithLabels(
		[]string{"exchange_rate"},
		float32(exchangeRate.MustFloat64()),
		[]metrics.Label{{Name: "denom", Value: denom}},
	)
}

// RecordAggregateExchangeRate records the median price gauge for a denom
func RecordMedianPrice(denom string, price sdk.Dec) {
	metrics.SetGaugeWithLabels(
		[]string{"median_price"},
		float32(price.MustFloat64()),
		[]metrics.Label{{Name: "denom", Value: denom}},
	)
}

// RecordAggregateExchangeRate records the median deviation price gauge for a denom
func RecordMedianDeviationPrice(denom string, price sdk.Dec) {
	metrics.SetGaugeWithLabels(
		[]string{"median_deviation_price"},
		float32(price.MustFloat64()),
		[]metrics.Label{{Name: "denom", Value: denom}},
	)
}
