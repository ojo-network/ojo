package metrics

import (
	"github.com/armon/go-metrics"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RecordMissCounter(operator sdk.ValAddress, missCounter uint64) {
	metrics.SetGaugeWithLabels(
		[]string{"miss_counter"},
		float32(missCounter),
		[]metrics.Label{{Name: "address", Value: operator.String()}},
	)
}

func RecordExchangeRate(denom string, exchangeRate sdk.Dec) {
	metrics.SetGaugeWithLabels(
		[]string{"exchange_rate"},
		float32(exchangeRate.MustFloat64()),
		[]metrics.Label{{Name: "denom", Value: denom}},
	)
}

func RecordMedianPrice(denom string, price sdk.Dec) {
	metrics.SetGaugeWithLabels(
		[]string{"median_price"},
		float32(price.MustFloat64()),
		[]metrics.Label{{Name: "denom", Value: denom}},
	)
}

func RecordMedianDeviationPrice(denom string, price sdk.Dec) {
	metrics.SetGaugeWithLabels(
		[]string{"median_deviation_price"},
		float32(price.MustFloat64()),
		[]metrics.Label{{Name: "denom", Value: denom}},
	)
}
