package oracle_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	ojoapp "github.com/ojo-network/ojo/app"
	appparams "github.com/ojo-network/ojo/app/params"
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

	ctx sdk.Context
	app *ojoapp.App
}

const (
	initialPower = int64(100)
)

func (s *IntegrationTestSuite) SetupTest() {
	require := s.Require()
	isCheckTx := false
	app := ojoapp.Setup(s.T())
	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{
		ChainID: fmt.Sprintf("test-chain-%s", tmrand.Str(4)),
	})

	oracle.InitGenesis(ctx, app.OracleKeeper, *types.DefaultGenesisState())

	sh := teststaking.NewHelper(s.T(), ctx, *app.StakingKeeper)
	sh.Denom = bondDenom

	// mint and send coins to validator
	require.NoError(app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initCoins))
	require.NoError(app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr1, initCoins))
	require.NoError(app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initCoins))
	require.NoError(app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr2, initCoins))

	// mint and send coins to oracle module to fill up reward pool
	require.NoError(app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initCoins))
	require.NoError(app.BankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, types.ModuleName, initCoins))

	sh.CreateValidatorWithValPower(valAddr1, valPubKey1, 70, true)
	sh.CreateValidatorWithValPower(valAddr2, valPubKey2, 30, true)

	staking.EndBlocker(ctx, *app.StakingKeeper)

	s.app = app
	s.ctx = ctx
}

// Test addresses
var (
	valPubKeys = simapp.CreateTestPubKeys(2)

	valPubKey1 = valPubKeys[0]
	pubKey1    = secp256k1.GenPrivKey().PubKey()
	addr1      = sdk.AccAddress(pubKey1.Address())
	valAddr1   = sdk.ValAddress(pubKey1.Address())

	valPubKey2 = valPubKeys[1]
	pubKey2    = secp256k1.GenPrivKey().PubKey()
	addr2      = sdk.AccAddress(pubKey2.Address())
	valAddr2   = sdk.ValAddress(pubKey2.Address())

	initTokens = sdk.TokensFromConsensusPower(initialPower, sdk.DefaultPowerReduction)
	initCoins  = sdk.NewCoins(sdk.NewCoin(bondDenom, initTokens))
)

func (s *IntegrationTestSuite) TestEndBlockerVoteThreshold() {
	app, ctx := s.app, s.ctx
	originalBlockHeight := ctx.BlockHeight()
	ctx = ctx.WithBlockHeight(1)
	preVoteBlockDiff := int64(app.OracleKeeper.VotePeriod(ctx) / 2)
	voteBlockDiff := int64(app.OracleKeeper.VotePeriod(ctx)/2 + 1)

	var (
		val1DecCoins sdk.DecCoins
		val2DecCoins sdk.DecCoins
		val1PreVotes types.AggregateExchangeRatePrevote
		val2PreVotes types.AggregateExchangeRatePrevote
		val1Votes    types.AggregateExchangeRateVote
		val2Votes    types.AggregateExchangeRateVote
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
	}

	val1PreVotes = types.AggregateExchangeRatePrevote{
		Hash:        "hash1",
		Voter:       valAddr1.String(),
		SubmitBlock: uint64(ctx.BlockHeight()),
	}
	val2PreVotes = types.AggregateExchangeRatePrevote{
		Hash:        "hash2",
		Voter:       valAddr2.String(),
		SubmitBlock: uint64(ctx.BlockHeight()),
	}

	val1Votes = types.AggregateExchangeRateVote{
		ExchangeRates: val1DecCoins,
		Voter:         valAddr1.String(),
	}
	val2Votes = types.AggregateExchangeRateVote{
		ExchangeRates: val2DecCoins,
		Voter:         valAddr2.String(),
	}

	// total voting power per denom is 100%
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr1, val1PreVotes)
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr2, val2PreVotes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + voteBlockDiff)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr1, val1Votes)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr2, val2Votes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	for _, denom := range app.OracleKeeper.AcceptList(ctx) {
		rate, err := app.OracleKeeper.GetExchangeRate(ctx, denom.SymbolDenom)
		s.Require().NoError(err)
		s.Require().Equal(sdk.MustNewDecFromStr("0.75"), rate)
	}

	// update prevotes' block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + preVoteBlockDiff)
	val1PreVotes.SubmitBlock = uint64(ctx.BlockHeight())
	val2PreVotes.SubmitBlock = uint64(ctx.BlockHeight())

	// total voting power per denom is 30%
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

	// update prevotes' block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + preVoteBlockDiff)
	val1PreVotes.SubmitBlock = uint64(ctx.BlockHeight())
	val2PreVotes.SubmitBlock = uint64(ctx.BlockHeight())

	// ojo has 100% power, and atom has 30%
	val1DecCoins = sdk.DecCoins{
		sdk.DecCoin{
			Denom:  "ojo",
			Amount: sdk.MustNewDecFromStr("1.0"),
		},
	}
	val2DecCoins = sdk.DecCoins{
		sdk.DecCoin{
			Denom:  "ojo",
			Amount: sdk.MustNewDecFromStr("0.5"),
		},
		sdk.DecCoin{
			Denom:  "atom",
			Amount: sdk.MustNewDecFromStr("0.5"),
		},
	}
	val1Votes.ExchangeRates = val1DecCoins
	val2Votes.ExchangeRates = val2DecCoins

	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr1, val1PreVotes)
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr2, val2PreVotes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + voteBlockDiff)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr1, val1Votes)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr2, val2Votes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	rate, err := app.OracleKeeper.GetExchangeRate(ctx, "ojo")
	s.Require().NoError(err)
	s.Require().Equal(sdk.MustNewDecFromStr("0.75"), rate)
	rate, err = app.OracleKeeper.GetExchangeRate(ctx, "atom")
	s.Require().ErrorIs(err, types.ErrUnknownDenom.Wrap("atom"))
	s.Require().Equal(sdk.ZeroDec(), rate)

	ctx = ctx.WithBlockHeight(originalBlockHeight)
}

func (s *IntegrationTestSuite) TestEndBlockerValidatorRewards() {
	app, ctx := s.app, s.ctx
	originalBlockHeight := ctx.BlockHeight()
	ctx = ctx.WithBlockHeight(1)
	preVoteBlockDiff := int64(app.OracleKeeper.VotePeriod(ctx) / 2)
	voteBlockDiff := int64(app.OracleKeeper.VotePeriod(ctx)/2 + 1)

	app.OracleKeeper.SetMandatoryList(ctx, types.DenomList{
		{
			BaseDenom:   bondDenom,
			SymbolDenom: appparams.DisplayDenom,
			Exponent:    uint32(6),
			RewardBand:  types.DefaultRewardBand,
		},
		{
			BaseDenom:   "ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9",
			SymbolDenom: "atom",
			Exponent:    uint32(6),
			RewardBand:  types.DefaultRewardBand,
		},
	})

	var (
		val1DecCoins sdk.DecCoins
		val2DecCoins sdk.DecCoins
		val1PreVotes types.AggregateExchangeRatePrevote
		val2PreVotes types.AggregateExchangeRatePrevote
		val1Votes    types.AggregateExchangeRateVote
		val2Votes    types.AggregateExchangeRateVote
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
	}

	val1PreVotes = types.AggregateExchangeRatePrevote{
		Hash:        "hash1",
		Voter:       valAddr1.String(),
		SubmitBlock: uint64(ctx.BlockHeight()),
	}
	val2PreVotes = types.AggregateExchangeRatePrevote{
		Hash:        "hash2",
		Voter:       valAddr2.String(),
		SubmitBlock: uint64(ctx.BlockHeight()),
	}

	val1Votes = types.AggregateExchangeRateVote{
		ExchangeRates: val1DecCoins,
		Voter:         valAddr1.String(),
	}
	val2Votes = types.AggregateExchangeRateVote{
		ExchangeRates: val2DecCoins,
		Voter:         valAddr2.String(),
	}

	// validator 1 and 2 vote on both currencies so both have 0 misses
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr1, val1PreVotes)
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr2, val2PreVotes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + voteBlockDiff)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr1, val1Votes)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr2, val2Votes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	s.Require().Equal(sdk.NewInt64DecCoin("uojo", 31), app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddr1).Rewards[0])
	s.Require().Equal(sdk.NewInt64DecCoin("uojo", 31), app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddr2).Rewards[0])

	// update prevotes' block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + preVoteBlockDiff)
	val1PreVotes.SubmitBlock = uint64(ctx.BlockHeight())
	val2PreVotes.SubmitBlock = uint64(ctx.BlockHeight())

	// validator 1 votes on both currencies to end up with 0 misses
	// validator 2 votes on 1 currency to end up with 1 misses
	val1DecCoins = sdk.DecCoins{
		sdk.DecCoin{
			Denom:  "ojo",
			Amount: sdk.MustNewDecFromStr("1.0"),
		},
		sdk.DecCoin{
			Denom:  "atom",
			Amount: sdk.MustNewDecFromStr("0.5"),
		},
	}
	val2DecCoins = sdk.DecCoins{
		sdk.DecCoin{
			Denom:  "ojo",
			Amount: sdk.MustNewDecFromStr("0.5"),
		},
	}
	val1Votes.ExchangeRates = val1DecCoins
	val2Votes.ExchangeRates = val2DecCoins

	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr1, val1PreVotes)
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr2, val2PreVotes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + voteBlockDiff)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr1, val1Votes)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr2, val2Votes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	s.Require().Equal(sdk.NewInt64DecCoin("uojo", 62), app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddr1).Rewards[0])
	s.Require().Equal(sdk.NewInt64DecCoin("uojo", 60), app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddr2).Rewards[0])

	// update prevotes' block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + preVoteBlockDiff)
	val1PreVotes.SubmitBlock = uint64(ctx.BlockHeight())
	val2PreVotes.SubmitBlock = uint64(ctx.BlockHeight())

	// validator 1 and 2 miss both currencies so validator 1 has 2 misses and
	// validator 2 has 3 misses
	val1Votes.ExchangeRates = sdk.DecCoins{}
	val2Votes.ExchangeRates = sdk.DecCoins{}

	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr1, val1PreVotes)
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr2, val2PreVotes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + voteBlockDiff)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr1, val1Votes)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr2, val2Votes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	s.Require().Equal(sdk.NewInt64DecCoin("uojo", 93), app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddr1).Rewards[0])
	s.Require().Equal(sdk.NewInt64DecCoin("uojo", 89), app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddr2).Rewards[0])

	ctx = ctx.WithBlockHeight(originalBlockHeight)
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
	blockHeight := ctx.BlockHeight()

	var historicStampPeriod int64 = 5
	var medianStampPeriod int64 = 20
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

func TestOracleTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
