package oracle_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	ojoapp "github.com/ojo-network/ojo/app"
	appparams "github.com/ojo-network/ojo/app/params"
	"github.com/ojo-network/ojo/tests/integration"
	"github.com/ojo-network/ojo/util/decmath"
	"github.com/ojo-network/ojo/x/oracle"
	"github.com/ojo-network/ojo/x/oracle/types"
)

const (
	displayDenom string = appparams.DisplayDenom
	bondDenom    string = appparams.BondDenom
)

type IntegrationTestSuite struct {
	suite.Suite

	ctx  sdk.Context
	app  *ojoapp.App
	keys []integration.TestValidatorKey
}

func (s *IntegrationTestSuite) SetupTest() {
	s.app, s.ctx, s.keys = integration.SetupAppWithContext(s.T())
	s.app.OracleKeeper.SetVoteThreshold(s.ctx, sdk.MustNewDecFromStr("0.4"))
}

func createVotes(hash string, val sdk.ValAddress, rates sdk.DecCoins, blockHeight uint64) (types.AggregateExchangeRatePrevote, types.AggregateExchangeRateVote) {
	preVote := types.AggregateExchangeRatePrevote{
		Hash:        hash,
		Voter:       val.String(),
		SubmitBlock: uint64(blockHeight),
	}
	vote := types.AggregateExchangeRateVote{
		ExchangeRates: rates,
		Voter:         val.String(),
	}
	return preVote, vote
}

func (s *IntegrationTestSuite) TestEndBlockerVoteThreshold() {
	app, ctx := s.app, s.ctx
	valAddr1, valAddr2, valAddr3 := s.keys[0].ValAddress, s.keys[1].ValAddress, s.keys[2].ValAddress
	ctx = ctx.WithBlockHeight(0)
	preVoteBlockDiff := int64(app.OracleKeeper.VotePeriod(ctx) / 2)
	voteBlockDiff := int64(app.OracleKeeper.VotePeriod(ctx)/2 + 1)

	var (
		val1DecCoins sdk.DecCoins
		val2DecCoins sdk.DecCoins
		val3DecCoins sdk.DecCoins
	)
	for _, denom := range app.OracleKeeper.AcceptList(ctx) {
		val1DecCoins = append(val1DecCoins, sdk.DecCoin{
			Denom:  denom.SymbolDenom,
			Amount: sdk.MustNewDecFromStr("1.0"),
		})
		val2DecCoins = append(val2DecCoins, sdk.DecCoin{
			Denom:  denom.SymbolDenom,
			Amount: sdk.MustNewDecFromStr("0.5"),
		})
		val3DecCoins = append(val3DecCoins, sdk.DecCoin{
			Denom:  denom.SymbolDenom,
			Amount: sdk.MustNewDecFromStr("0.6"),
		})
	}

	// add junk coin and ensure ballot still is counted
	junkCoin := sdk.DecCoin{
		Denom:  "JUNK",
		Amount: sdk.MustNewDecFromStr("0.05"),
	}
	val1DecCoins = append(val1DecCoins, junkCoin)
	val2DecCoins = append(val2DecCoins, junkCoin)
	val3DecCoins = append(val3DecCoins, junkCoin)

	h := uint64(ctx.BlockHeight())
	val1PreVotes, val1Votes := createVotes("hash1", valAddr1, val1DecCoins, h)
	val2PreVotes, val2Votes := createVotes("hash2", valAddr2, val2DecCoins, h)
	val3PreVotes, val3Votes := createVotes("hash3", valAddr3, val3DecCoins, h)

	// total voting power per denom is 100%
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr1, val1PreVotes)
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr2, val2PreVotes)
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr3, val3PreVotes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + voteBlockDiff)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr1, val1Votes)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr2, val2Votes)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr3, val3Votes)
	err := oracle.EndBlocker(ctx, app.OracleKeeper)
	s.Require().NoError(err)

	for _, denom := range app.OracleKeeper.AcceptList(ctx) {
		rate, err := app.OracleKeeper.GetExchangeRate(ctx, denom.SymbolDenom)
		s.Require().NoError(err)
		s.Require().Equal(sdk.MustNewDecFromStr("1.0"), rate)
	}

	// Test: only val2 votes (has 39% vote power).
	// Total voting power per denom must be bigger or equal than 40% (see SetupTest).
	// So if only val2 votes, we won't have any prices next block.
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + preVoteBlockDiff)
	h = uint64(ctx.BlockHeight())
	val2PreVotes.SubmitBlock = h

	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr2, val2PreVotes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + voteBlockDiff)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr2, val2Votes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	for _, denom := range app.OracleKeeper.AcceptList(ctx) {
		rate, err := app.OracleKeeper.GetExchangeRate(ctx, denom.SymbolDenom)
		s.Require().ErrorIs(err, types.ErrUnknownDenom.Wrap(denom.SymbolDenom))
		s.Require().Equal(sdk.ZeroDec(), rate)
	}

	// Test: val2 and val3 votes.
	// now we will have 40% of the power, so now we should have prices
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + preVoteBlockDiff)
	h = uint64(ctx.BlockHeight())
	val2PreVotes.SubmitBlock = h
	val3PreVotes.SubmitBlock = h

	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr2, val2PreVotes)
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr3, val3PreVotes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + voteBlockDiff)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr2, val2Votes)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr3, val3Votes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	for _, denom := range app.OracleKeeper.AcceptList(ctx) {
		rate, err := app.OracleKeeper.GetExchangeRate(ctx, denom.SymbolDenom)
		s.Require().NoError(err)
		s.Require().Equal(sdk.MustNewDecFromStr("0.5"), rate)
	}

	// Test: val1 and val2 vote again
	// umee has 69.9% power, and atom has 30%, so we should have price for umee, but not for atom
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + preVoteBlockDiff)
	h = uint64(ctx.BlockHeight())
	val1PreVotes.SubmitBlock = h
	val2PreVotes.SubmitBlock = h

	val1Votes.ExchangeRates = sdk.DecCoins{
		sdk.NewDecCoinFromDec("ojo", sdk.MustNewDecFromStr("1.0")),
	}
	val2Votes.ExchangeRates = sdk.DecCoins{
		sdk.NewDecCoinFromDec("atom", sdk.MustNewDecFromStr("0.5")),
	}

	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr1, val1PreVotes)
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr2, val2PreVotes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + voteBlockDiff)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr1, val1Votes)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr2, val2Votes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	rate, err := app.OracleKeeper.GetExchangeRate(ctx, "ojo")
	s.Require().NoError(err)
	s.Require().Equal(sdk.MustNewDecFromStr("1.0"), rate)
	rate, err = app.OracleKeeper.GetExchangeRate(ctx, "atom")
	s.Require().ErrorIs(err, types.ErrUnknownDenom.Wrap("atom"))
	s.Require().Equal(sdk.ZeroDec(), rate)
}

func (s *IntegrationTestSuite) TestEndBlockerValidatorRewards() {
	app, ctx := s.app, s.ctx
	valAddr1, valAddr2, valAddr3 := s.keys[0].ValAddress, s.keys[1].ValAddress, s.keys[2].ValAddress
	preVoteBlockDiff := int64(app.OracleKeeper.VotePeriod(ctx) / 2)
	voteBlockDiff := int64(app.OracleKeeper.VotePeriod(ctx)/2 + 1)

	// start test in new slash window
	ctx = ctx.WithBlockHeight(int64(app.OracleKeeper.SlashWindow(ctx)))
	oracle.EndBlocker(ctx, app.OracleKeeper)

	app.OracleKeeper.SetMandatoryList(ctx, types.DenomList{
		{
			BaseDenom:   bondDenom,
			SymbolDenom: appparams.DisplayDenom,
			Exponent:    uint32(6),
		},
		{
			BaseDenom:   "ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9",
			SymbolDenom: "atom",
			Exponent:    uint32(6),
		},
	})

	var (
		val1DecCoins sdk.DecCoins
		val2DecCoins sdk.DecCoins
		val3DecCoins sdk.DecCoins
	)
	for _, denom := range app.OracleKeeper.AcceptList(ctx) {
		val1DecCoins = append(val1DecCoins, sdk.DecCoin{
			Denom:  denom.SymbolDenom,
			Amount: sdk.MustNewDecFromStr("0.6"),
		})
		val2DecCoins = append(val2DecCoins, sdk.DecCoin{
			Denom:  denom.SymbolDenom,
			Amount: sdk.MustNewDecFromStr("0.6"),
		})
		val3DecCoins = append(val3DecCoins, sdk.DecCoin{
			Denom:  denom.SymbolDenom,
			Amount: sdk.MustNewDecFromStr("0.6"),
		})
	}

	h := uint64(ctx.BlockHeight())
	val1PreVotes, val1Votes := createVotes("hash1", valAddr1, val1DecCoins, h)
	val2PreVotes, val2Votes := createVotes("hash2", valAddr2, val2DecCoins, h)
	val3PreVotes, val3Votes := createVotes("hash3", valAddr3, val3DecCoins, h)
	// validator 1, 2, and 3 vote on both currencies so all have 0 misses
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr1, val1PreVotes)
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr2, val2PreVotes)
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr3, val3PreVotes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + voteBlockDiff)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr1, val1Votes)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr2, val2Votes)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr3, val3Votes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	s.Require().Equal(sdk.NewInt64DecCoin("uojo", 142), app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddr1).Rewards[0])
	s.Require().Equal(sdk.NewInt64DecCoin("uojo", 142), app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddr2).Rewards[0])
	s.Require().Equal(sdk.NewInt64DecCoin("uojo", 142), app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddr3).Rewards[0])

	// update prevotes' block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + preVoteBlockDiff)
	val1PreVotes.SubmitBlock = uint64(ctx.BlockHeight())
	val2PreVotes.SubmitBlock = uint64(ctx.BlockHeight())
	val3PreVotes.SubmitBlock = uint64(ctx.BlockHeight())

	// validator 1 and 3 votes on both currencies to end up with 0 misses
	// validator 2 votes on 1 currency to end up with 1 misses
	val1DecCoins = sdk.DecCoins{
		sdk.DecCoin{
			Denom:  "ojo",
			Amount: sdk.MustNewDecFromStr("0.6"),
		},
		sdk.DecCoin{
			Denom:  "atom",
			Amount: sdk.MustNewDecFromStr("0.6"),
		},
	}
	val2DecCoins = sdk.DecCoins{
		sdk.DecCoin{
			Denom:  "ojo",
			Amount: sdk.MustNewDecFromStr("0.6"),
		},
	}
	val3DecCoins = sdk.DecCoins{
		sdk.DecCoin{
			Denom:  "ojo",
			Amount: sdk.MustNewDecFromStr("0.6"),
		},
		sdk.DecCoin{
			Denom:  "atom",
			Amount: sdk.MustNewDecFromStr("0.6"),
		},
	}
	val1Votes.ExchangeRates = val1DecCoins
	val2Votes.ExchangeRates = val2DecCoins
	val3Votes.ExchangeRates = val3DecCoins

	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr1, val1PreVotes)
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr2, val2PreVotes)
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr3, val3PreVotes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + voteBlockDiff)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr1, val1Votes)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr2, val2Votes)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr3, val3Votes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	s.Require().Equal(sdk.NewInt64DecCoin("uojo", 284), app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddr1).Rewards[0])
	s.Require().Equal(sdk.NewInt64DecCoin("uojo", 275), app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddr2).Rewards[0])
	s.Require().Equal(sdk.NewInt64DecCoin("uojo", 284), app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddr3).Rewards[0])

	// update prevotes' block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + preVoteBlockDiff)
	val1PreVotes.SubmitBlock = uint64(ctx.BlockHeight())
	val2PreVotes.SubmitBlock = uint64(ctx.BlockHeight())

	// validator 1, 2, and 3 miss both currencies so validator 1 and 3 has 2 misses and
	// validator 2 has 3 misses
	val1Votes.ExchangeRates = sdk.DecCoins{}
	val2Votes.ExchangeRates = sdk.DecCoins{}
	val3Votes.ExchangeRates = sdk.DecCoins{}

	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr1, val1PreVotes)
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr2, val2PreVotes)
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr3, val3PreVotes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + voteBlockDiff)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr1, val1Votes)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr2, val2Votes)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr3, val3Votes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	s.Require().Equal(sdk.NewInt64DecCoin("uojo", 426), app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddr1).Rewards[0])
	s.Require().Equal(sdk.NewInt64DecCoin("uojo", 408), app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddr2).Rewards[0])
	s.Require().Equal(sdk.NewInt64DecCoin("uojo", 426), app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddr3).Rewards[0])
}

var exchangeRates = map[string][]sdk.Dec{
	"ATOM": {
		sdk.MustNewDecFromStr("12.99"),
		sdk.MustNewDecFromStr("12.22"),
		sdk.MustNewDecFromStr("13.1"),
		sdk.MustNewDecFromStr("11.6"),
	},
	"OJO": {
		sdk.MustNewDecFromStr("1.89"),
		sdk.MustNewDecFromStr("2.05"),
		sdk.MustNewDecFromStr("2.34"),
		sdk.MustNewDecFromStr("1.71"),
	},
}

func (s *IntegrationTestSuite) TestEndblockerHistoracle() {
	app, ctx := s.app, s.ctx
	valAddr1 := s.keys[0].ValAddress
	blockHeight := ctx.BlockHeight()

	var historicStampPeriod int64 = 3
	var medianStampPeriod int64 = 12
	var maximumPriceStamps int64 = 4
	var maximumMedianStamps int64 = 3

	app.OracleKeeper.SetHistoricStampPeriod(ctx, uint64(historicStampPeriod))
	app.OracleKeeper.SetMedianStampPeriod(ctx, uint64(medianStampPeriod))
	app.OracleKeeper.SetMaximumPriceStamps(ctx, uint64(maximumPriceStamps))
	app.OracleKeeper.SetMaximumMedianStamps(ctx, uint64(maximumMedianStamps))

	// Start at the last block of the first stamp period
	blockHeight += medianStampPeriod
	blockHeight += -1
	ctx = ctx.WithBlockHeight(blockHeight)

	for i := int64(0); i <= maximumMedianStamps; i++ {
		for j := int64(0); j < maximumPriceStamps; j++ {

			blockHeight += historicStampPeriod
			ctx = ctx.WithBlockHeight(blockHeight)

			decCoins := sdk.DecCoins{}
			for denom, prices := range exchangeRates {
				decCoins = append(decCoins, sdk.DecCoin{
					Denom:  denom,
					Amount: prices[j],
				})
			}

			vote := types.AggregateExchangeRateVote{
				ExchangeRates: decCoins,
				Voter:         valAddr1.String(),
			}
			app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr1, vote)
			oracle.EndBlocker(ctx, app.OracleKeeper)
		}

		for denom, denomRates := range exchangeRates {
			// check median
			expectedMedian, err := decmath.Median(denomRates)
			s.Require().NoError(err)

			medians := app.OracleKeeper.AllMedianPrices(ctx)
			medians = *medians.FilterByBlock(uint64(blockHeight)).FilterByDenom(denom)
			actualMedian := medians[0].ExchangeRate.Amount
			s.Require().Equal(expectedMedian, actualMedian)

			// check median deviation
			expectedMedianDeviation, err := decmath.MedianDeviation(actualMedian, denomRates)
			s.Require().NoError(err)

			medianDeviations := app.OracleKeeper.AllMedianDeviationPrices(ctx)
			medianDeviations = *medianDeviations.FilterByBlock(uint64(blockHeight)).FilterByDenom(denom)
			actualMedianDeviation := medianDeviations[0].ExchangeRate.Amount
			s.Require().Equal(expectedMedianDeviation, actualMedianDeviation)
		}
	}
	numberOfAssets := int64(len(exchangeRates))

	historicPrices := app.OracleKeeper.AllHistoricPrices(ctx)
	s.Require().Equal(maximumPriceStamps*numberOfAssets, int64(len(historicPrices)))

	for i := int64(0); i < maximumPriceStamps; i++ {
		expectedBlockNum := blockHeight - (historicStampPeriod * (maximumPriceStamps - int64(i+1)))
		actualBlockNum := historicPrices[i].BlockNum
		s.Require().Equal(expectedBlockNum, int64(actualBlockNum))
	}

	medians := app.OracleKeeper.AllMedianPrices(ctx)
	s.Require().Equal(maximumMedianStamps*numberOfAssets, int64(len(medians)))

	for i := int64(0); i < maximumMedianStamps; i++ {
		expectedBlockNum := blockHeight - (medianStampPeriod * (maximumMedianStamps - int64(i+1)))
		actualBlockNum := medians[i].BlockNum
		s.Require().Equal(expectedBlockNum, int64(actualBlockNum))
	}

	medianDeviations := app.OracleKeeper.AllMedianPrices(ctx)
	s.Require().Equal(maximumMedianStamps*numberOfAssets, int64(len(medianDeviations)))

	for i := int64(0); i < maximumMedianStamps; i++ {
		expectedBlockNum := blockHeight - (medianStampPeriod * (maximumMedianStamps - int64(i+1)))
		actualBlockNum := medianDeviations[i].BlockNum
		s.Require().Equal(expectedBlockNum, int64(actualBlockNum))
	}
}

func (s *IntegrationTestSuite) TestUpdateOracleParams() {
	app, ctx := s.app, s.ctx
	blockHeight := ctx.BlockHeight()

	// Schedule param update plan for current block height
	err := app.OracleKeeper.ScheduleParamUpdatePlan(
		ctx,
		types.ParamUpdatePlan{
			Keys:   []string{"VoteThreshold"},
			Height: blockHeight,
			Changes: types.Params{
				VoteThreshold: sdk.NewDecWithPrec(40, 2),
			},
		},
	)
	s.Require().NoError(err)
	_, found := s.app.OracleKeeper.GetParamUpdatePlan(s.ctx)
	s.Require().Equal(true, found)

	// Check Vote Threshold was updated
	oracle.EndBlocker(ctx, app.OracleKeeper)
	s.Require().Equal(sdk.NewDecWithPrec(40, 2), app.OracleKeeper.VoteThreshold(ctx))

	// Schedule param update plan for current block height and then cancel it
	err = app.OracleKeeper.ScheduleParamUpdatePlan(
		ctx,
		types.ParamUpdatePlan{
			Keys:   []string{"VoteThreshold"},
			Height: blockHeight,
			Changes: types.Params{
				VoteThreshold: sdk.NewDecWithPrec(50, 2),
			},
		},
	)
	s.Require().NoError(err)
	_, found = s.app.OracleKeeper.GetParamUpdatePlan(s.ctx)
	s.Require().Equal(true, found)

	// Cancel update
	err = app.OracleKeeper.ClearParamUpdatePlan(ctx)
	s.Require().NoError(err)
	_, found = s.app.OracleKeeper.GetParamUpdatePlan(s.ctx)
	s.Require().Equal(false, found)

	// Check Vote Threshold wasn't updated
	oracle.EndBlocker(ctx, app.OracleKeeper)
	s.Require().Equal(sdk.NewDecWithPrec(40, 2), app.OracleKeeper.VoteThreshold(ctx))
}

func TestOracleTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
