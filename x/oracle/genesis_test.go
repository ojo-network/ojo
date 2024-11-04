package oracle_test

import (
	"testing"

	"gotest.tools/v3/assert"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	appparams "github.com/ojo-network/ojo/app/params"
	"github.com/ojo-network/ojo/tests/integration"
	"github.com/stretchr/testify/suite"

	ojoapp "github.com/ojo-network/ojo/app"
	"github.com/ojo-network/ojo/x/oracle"
	"github.com/ojo-network/ojo/x/oracle/types"
)

const (
	valoperAddr = "ojovaloper18wvrhmvm8av6s34q48yees3yu6xglwhmr4tzdn"
	addr        = "ojo1ner6kc63xl903wrv2n8p9mtun79gegjld93lx0"
	denom       = "ojo"
)

type IntegrationTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *ojoapp.App
}

var exchangeRate = math.LegacyMustNewDecFromStr("8.8")

func (s *IntegrationTestSuite) SetupTest() {
	s.app, s.ctx, _ = integration.SetupAppWithContext(s.T())
}

func (s *IntegrationTestSuite) TestGenesis_InitGenesis() {
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
					{
						ExchangeRate: &sdk.DecCoin{
							Denom:  denom,
							Amount: exchangeRate,
						},
						BlockNum: 0,
					},
				},
				Medians: []types.PriceStamp{
					{
						ExchangeRate: &sdk.DecCoin{
							Denom:  denom,
							Amount: exchangeRate,
						},
						BlockNum: 0,
					},
				},
				MedianDeviations: []types.PriceStamp{
					{
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

func (s *IntegrationTestSuite) TestGenesis_ExportGenesis() {
	keeper, ctx := s.app.OracleKeeper, s.ctx

	params := types.DefaultParams()

	feederDelegations := []types.FeederDelegation{
		{
			FeederAddress:    addr,
			ValidatorAddress: valoperAddr,
		},
	}
	exchangeRates := sdk.DecCoins{
		{
			Denom:  appparams.DisplayDenom,
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
		ExchangeRates:                 exchangeRates,
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
	assert.DeepEqual(s.T(), exchangeRates, result.ExchangeRates)
	assert.DeepEqual(s.T(), missCounters[0], result.MissCounters[0])
	assert.DeepEqual(s.T(), aggregateExchangeRatePrevotes, result.AggregateExchangeRatePrevotes)
	assert.DeepEqual(s.T(), aggregateExchangeRateVotes, result.AggregateExchangeRateVotes)
	assert.DeepEqual(s.T(), medians, result.Medians)
	assert.DeepEqual(s.T(), historicPrices, result.HistoricPrices)
	assert.DeepEqual(s.T(), medianDeviations, result.MedianDeviations)
}

func TestOracleTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
