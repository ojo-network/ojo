package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(
	params Params,
	rates sdk.DecCoins,
	feederDelegations []FeederDelegation,
	missCounters []MissCounter,
	aggregateExchangeRatePrevotes []AggregateExchangeRatePrevote,
	aggregateExchangeRateVotes []AggregateExchangeRateVote,
	historicPrices []PriceStamp,
	medianPrices []PriceStamp,
	medianDeviationPrices []PriceStamp,
) *GenesisState {
	return &GenesisState{
		Params:                        params,
		ExchangeRates:                 rates,
		FeederDelegations:             feederDelegations,
		MissCounters:                  missCounters,
		AggregateExchangeRatePrevotes: aggregateExchangeRatePrevotes,
		AggregateExchangeRateVotes:    aggregateExchangeRateVotes,
		HistoricPrices:                historicPrices,
		Medians:                       medianPrices,
		MedianDeviations:              medianDeviationPrices,
	}
}

// DefaultGenesisState returns the default genesesis state for the x/oracle
// module.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:                        DefaultParams(),
		ExchangeRates:                 sdk.DecCoins{},
		FeederDelegations:             []FeederDelegation{},
		MissCounters:                  []MissCounter{},
		AggregateExchangeRatePrevotes: []AggregateExchangeRatePrevote{},
		AggregateExchangeRateVotes:    []AggregateExchangeRateVote{},
		HistoricPrices:                []PriceStamp{},
		Medians:                       []PriceStamp{},
		MedianDeviations:              []PriceStamp{},
	}
}

// ValidateGenesis validates the oracle genesis state.
func ValidateGenesis(data *GenesisState) error {
	return data.Params.Validate()
}

// GetGenesisStateFromAppState returns x/oracle GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc codec.JSONCodec, appState map[string]json.RawMessage) *GenesisState {
	var genesisState GenesisState

	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}

	return &genesisState
}
