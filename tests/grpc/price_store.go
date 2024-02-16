package grpc

import (
	"fmt"

	"cosmossdk.io/math"
	"github.com/ojo-network/ojo/util/decmath"
)

// PriceStore is the in memory store for prices
// recorded by the price_listener with helper methods
// for calculating and verifying medians and median deviations
type PriceStore struct {
	historicStamps   map[string][]math.LegacyDec
	medians          map[string]math.LegacyDec
	medianDeviations map[string]math.LegacyDec
}

func NewPriceStore() *PriceStore {
	return &PriceStore{
		historicStamps:   map[string][]math.LegacyDec{},
		medians:          map[string]math.LegacyDec{},
		medianDeviations: map[string]math.LegacyDec{},
	}
}

func (ps *PriceStore) addStamp(denom string, stamp math.LegacyDec) {
	if _, ok := ps.historicStamps[denom]; !ok {
		ps.historicStamps[denom] = []math.LegacyDec{}
	}
	ps.historicStamps[denom] = append(ps.historicStamps[denom], stamp)
}

func (ps *PriceStore) checkMedians() error {
	for denom, stamps := range ps.historicStamps {
		calcMedian, err := decmath.Median(stamps)
		if err != nil {
			return err
		}
		if !ps.medians[denom].Equal(calcMedian) {
			return fmt.Errorf(
				"expected %d for the %s median but got %d",
				ps.medians[denom],
				denom,
				calcMedian,
			)
		}
	}
	return nil
}

func (ps *PriceStore) checkMedianDeviations() error {
	for denom, median := range ps.medians {
		calcMedianDeviation, err := decmath.MedianDeviation(median, ps.historicStamps[denom])
		if err != nil {
			return err
		}
		if !ps.medianDeviations[denom].Equal(calcMedianDeviation) {
			return fmt.Errorf(
				"expected %d for the %s median deviation but got %d",
				ps.medianDeviations[denom],
				denom,
				calcMedianDeviation,
			)
		}
	}
	return nil
}
