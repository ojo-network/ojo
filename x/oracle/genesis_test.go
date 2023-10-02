package oracle_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"gotest.tools/assert"

	"github.com/ojo-network/ojo/x/oracle"
	"github.com/ojo-network/ojo/x/oracle/types"
)

func (s *GenesisTestSuite) TestGenesis_InitGenesis() {
	fmt.Println("init genesis test")
	keeper, ctx := s.app.OracleKeeper, s.ctx

	tcs := []struct {
		name      string
		g         types.GenesisState
		expectErr bool
		errMsg    string
	}{
		{
			"FeederDelegations.FeederAddress: empty address",
			types.GenesisState{
				FeederDelegations: []types.FeederDelegation{
					{
						FeederAddress: "",
					},
				},
			},
			true,
			"empty address string is not allowed",
		},
		{
			"FeederDelegations.ValidatorAddress: empty address",
			types.GenesisState{
				FeederDelegations: []types.FeederDelegation{
					{
						FeederAddress:    addr,
						ValidatorAddress: "",
					},
				},
			},
			true,
			"empty address string is not allowed",
		},
		{
			"valid",
			types.GenesisState{
				Params: types.DefaultParams(),
				ExchangeRates: []sdk.DecCoin{
					{
						Denom:  denom,
						Amount: exchangeRate,
					},
				},
				HistoricPrices: []types.PriceStamp{
					types.PriceStamp{
						ExchangeRate: &sdk.DecCoin{
							Denom:  denom,
							Amount: exchangeRate,
						},
						BlockNum: 0,
					},
				},
				Medians: []types.PriceStamp{
					types.PriceStamp{
						ExchangeRate: &sdk.DecCoin{
							Denom:  denom,
							Amount: exchangeRate,
						},
						BlockNum: 0,
					},
				},
				MedianDeviations: []types.PriceStamp{
					types.PriceStamp{
						ExchangeRate: &sdk.DecCoin{
							Denom:  denom,
							Amount: exchangeRate,
						},
						BlockNum: 0,
					},
				},
			},
			false,
			"",
		},
		{
			"FeederDelegations.ValidatorAddress: empty address",
			types.GenesisState{
				MissCounters: []types.MissCounter{
					{
						ValidatorAddress: "",
					},
				},
			},
			true,
			"empty address string is not allowed",
		},
		{
			"AggregateExchangeRatePrevotes.Voter: empty address",
			types.GenesisState{
				AggregateExchangeRatePrevotes: []types.AggregateExchangeRatePrevote{
					{
						Voter: "",
					},
				},
			},
			true,
			"empty address string is not allowed",
		},
		{
			"AggregateExchangeRateVotes.Voter: empty address",
			types.GenesisState{
				AggregateExchangeRateVotes: []types.AggregateExchangeRateVote{
					{
						Voter: "",
					},
				},
			},
			true,
			"empty address string is not allowed",
		},
	}

	for _, tc := range tcs {
		s.Run(
			tc.name, func() {
				if tc.expectErr {
					s.Assertions.PanicsWithError(tc.errMsg, func() { oracle.InitGenesis(ctx, keeper, tc.g) })
				} else {
					s.Assertions.NotPanics(func() { oracle.InitGenesis(ctx, keeper, tc.g) })
				}
			},
		)
	}
}

func (s *GenesisTestSuite) TestGenesis_ExportGenesis() {
	fmt.Println("export genesis test")
	keeper, ctx := s.app.OracleKeeper, s.ctx

	params := types.DefaultParams()

	feederDelegations := []types.FeederDelegation{
		{
			FeederAddress:    addr,
			ValidatorAddress: valoperAddr,
		},
	}
	exchangeRateTuples := []sdk.DecCoin{
		{
			Denom:  upperDenom,
			Amount: exchangeRate,
		},
	}
	missCounters := []types.MissCounter{
		{
			ValidatorAddress: valoperAddr,
		},
	}
	aggregateExchangeRatePrevotes := []types.AggregateExchangeRatePrevote{
		{
			Voter: valoperAddr,
		},
	}
	aggregateExchangeRateVotes := []types.AggregateExchangeRateVote{
		{
			Voter: valoperAddr,
		},
	}

	historicPrices := []types.PriceStamp{
		{
			ExchangeRate: &sdk.DecCoin{
				Denom:  denom,
				Amount: exchangeRate,
			},
			BlockNum: 0,
		},
	}

	medians := []types.PriceStamp{
		{
			ExchangeRate: &sdk.DecCoin{
				Denom:  denom,
				Amount: exchangeRate,
			},
			BlockNum: 0,
		},
	}

	medianDeviations := []types.PriceStamp{
		{
			ExchangeRate: &sdk.DecCoin{
				Denom:  denom,
				Amount: exchangeRate,
			},
			BlockNum: 0,
		},
	}

	genesisState := types.GenesisState{
		Params:                        params,
		FeederDelegations:             feederDelegations,
		ExchangeRates:                 exchangeRateTuples,
		MissCounters:                  missCounters,
		AggregateExchangeRatePrevotes: aggregateExchangeRatePrevotes,
		AggregateExchangeRateVotes:    aggregateExchangeRateVotes,
		Medians:                       medians,
		HistoricPrices:                historicPrices,
		MedianDeviations:              medianDeviations,
	}

	oracle.InitGenesis(ctx, keeper, genesisState)

	result := oracle.ExportGenesis(s.ctx, s.app.OracleKeeper)
	assert.DeepEqual(s.T(), params, result.Params)
	assert.DeepEqual(s.T(), feederDelegations, result.FeederDelegations)
	assert.DeepEqual(s.T(), exchangeRateTuples, result.ExchangeRates)
	assert.DeepEqual(s.T(), missCounters, result.MissCounters)
	assert.DeepEqual(s.T(), aggregateExchangeRatePrevotes, result.AggregateExchangeRatePrevotes)
	assert.DeepEqual(s.T(), aggregateExchangeRateVotes, result.AggregateExchangeRateVotes)
	assert.DeepEqual(s.T(), medians, result.Medians)
	assert.DeepEqual(s.T(), historicPrices, result.HistoricPrices)
	assert.DeepEqual(s.T(), medianDeviations, result.MedianDeviations)
	assert.DeepEqual(s.T(), hacp, result.AvgCounterParams)
}
