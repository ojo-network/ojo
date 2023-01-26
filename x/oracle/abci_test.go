package oracle_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	ojoapp "github.com/ojo-network/ojo/app"
	appparams "github.com/ojo-network/ojo/app/params"
	"github.com/ojo-network/ojo/x/oracle"
	"github.com/ojo-network/ojo/x/oracle/keeper"
	"github.com/ojo-network/ojo/x/oracle/types"
)

const (
	bondDenom        string = appparams.BondDenom
	preVoteBlockDiff int64  = 2
	voteBlockDiff    int64  = 3
)

type IntegrationTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *ojoapp.App
}

const (
	initialPower = int64(100)
)

// clearHistoricPrices deletes all historic prices of a given denom in the store.
func clearHistoricPrices(
	ctx sdk.Context,
	k keeper.Keeper,
	denom string,
) {
	stampPeriod := int(k.HistoricStampPeriod(ctx))
	numStamps := int(k.MaximumPriceStamps(ctx))
	for i := 0; i < numStamps; i++ {
		k.DeleteHistoricPrice(ctx, denom, uint64(ctx.BlockHeight())-uint64(i*stampPeriod))
	}
}

// clearHistoricMedians deletes all historic medians of a given denom in the store.
func clearHistoricMedians(
	ctx sdk.Context,
	k keeper.Keeper,
	denom string,
) {
	stampPeriod := int(k.MedianStampPeriod(ctx))
	numStamps := int(k.MaximumMedianStamps(ctx))
	for i := 0; i < numStamps; i++ {
		k.DeleteHistoricMedian(ctx, denom, uint64(ctx.BlockHeight())-uint64(i*stampPeriod))
	}
}

// clearHistoricMedianDeviations deletes all historic median deviations of a given
// denom in the store.
func clearHistoricMedianDeviations(
	ctx sdk.Context,
	k keeper.Keeper,
	denom string,
) {
	stampPeriod := int(k.MedianStampPeriod(ctx))
	numStamps := int(k.MaximumMedianStamps(ctx))
	for i := 0; i < numStamps; i++ {
		k.DeleteHistoricMedianDeviation(ctx, denom, uint64(ctx.BlockHeight())-uint64(i*stampPeriod))
	}
}

// SetupTest will create and supply two validators with %100
// of the consensus power worth of tokens split 70/30.
func (s *IntegrationTestSuite) SetupTest() {
	require := s.Require()
	isCheckTx := false
	app := ojoapp.Setup(s.T())
	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{
		ChainID: fmt.Sprintf("test-chain-%s", tmrand.Str(4)),
		Height:  6,
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

	var (
		val1Tuples   types.ExchangeRateTuples
		val2Tuples   types.ExchangeRateTuples
		val1PreVotes types.AggregateExchangeRatePrevote
		val2PreVotes types.AggregateExchangeRatePrevote
		val1Votes    types.AggregateExchangeRateVote
		val2Votes    types.AggregateExchangeRateVote
	)
	for _, denom := range app.OracleKeeper.AcceptList(ctx) {
		val1Tuples = append(val1Tuples, types.ExchangeRateTuple{
			Denom:        denom.SymbolDenom,
			ExchangeRate: sdk.MustNewDecFromStr("1.0"),
		})
		val2Tuples = append(val2Tuples, types.ExchangeRateTuple{
			Denom:        denom.SymbolDenom,
			ExchangeRate: sdk.MustNewDecFromStr("0.5"),
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
		ExchangeRateTuples: val1Tuples,
		Voter:              valAddr1.String(),
	}
	val2Votes = types.AggregateExchangeRateVote{
		ExchangeRateTuples: val2Tuples,
		Voter:              valAddr2.String(),
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
		s.Require().ErrorIs(err, sdkerrors.Wrap(types.ErrUnknownDenom, denom.SymbolDenom))
		s.Require().Equal(sdk.ZeroDec(), rate)
	}

	// update prevotes' block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + preVoteBlockDiff)
	val1PreVotes.SubmitBlock = uint64(ctx.BlockHeight())
	val2PreVotes.SubmitBlock = uint64(ctx.BlockHeight())

	// ojo has 100% power, and atom has 30%
	val1Tuples = types.ExchangeRateTuples{
		types.ExchangeRateTuple{
			Denom:        "ojo",
			ExchangeRate: sdk.MustNewDecFromStr("1.0"),
		},
	}
	val2Tuples = types.ExchangeRateTuples{
		types.ExchangeRateTuple{
			Denom:        "ojo",
			ExchangeRate: sdk.MustNewDecFromStr("0.5"),
		},
		types.ExchangeRateTuple{
			Denom:        "atom",
			ExchangeRate: sdk.MustNewDecFromStr("0.5"),
		},
	}
	val1Votes.ExchangeRateTuples = val1Tuples
	val2Votes.ExchangeRateTuples = val2Tuples

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
	s.Require().ErrorIs(err, sdkerrors.Wrap(types.ErrUnknownDenom, "atom"))
	s.Require().Equal(sdk.ZeroDec(), rate)
}

func (s *IntegrationTestSuite) TestEndBlockerValidatorRewards() {
	app, ctx := s.app, s.ctx

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
		val1Tuples   types.ExchangeRateTuples
		val2Tuples   types.ExchangeRateTuples
		val1PreVotes types.AggregateExchangeRatePrevote
		val2PreVotes types.AggregateExchangeRatePrevote
		val1Votes    types.AggregateExchangeRateVote
		val2Votes    types.AggregateExchangeRateVote
	)
	for _, denom := range app.OracleKeeper.AcceptList(ctx) {
		val1Tuples = append(val1Tuples, types.ExchangeRateTuple{
			Denom:        denom.SymbolDenom,
			ExchangeRate: sdk.MustNewDecFromStr("1.0"),
		})
		val2Tuples = append(val2Tuples, types.ExchangeRateTuple{
			Denom:        denom.SymbolDenom,
			ExchangeRate: sdk.MustNewDecFromStr("0.5"),
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
		ExchangeRateTuples: val1Tuples,
		Voter:              valAddr1.String(),
	}
	val2Votes = types.AggregateExchangeRateVote{
		ExchangeRateTuples: val2Tuples,
		Voter:              valAddr2.String(),
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
	val1Tuples = types.ExchangeRateTuples{
		types.ExchangeRateTuple{
			Denom:        "ojo",
			ExchangeRate: sdk.MustNewDecFromStr("1.0"),
		},
		types.ExchangeRateTuple{
			Denom:        "atom",
			ExchangeRate: sdk.MustNewDecFromStr("0.5"),
		},
	}
	val2Tuples = types.ExchangeRateTuples{
		types.ExchangeRateTuple{
			Denom:        "ojo",
			ExchangeRate: sdk.MustNewDecFromStr("0.5"),
		},
	}
	val1Votes.ExchangeRateTuples = val1Tuples
	val2Votes.ExchangeRateTuples = val2Tuples

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
	val1Votes.ExchangeRateTuples = types.ExchangeRateTuples{}
	val2Votes.ExchangeRateTuples = types.ExchangeRateTuples{}

	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr1, val1PreVotes)
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr2, val2PreVotes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + voteBlockDiff)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr1, val1Votes)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr2, val2Votes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	s.Require().Equal(sdk.NewInt64DecCoin("uojo", 93), app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddr1).Rewards[0])
	s.Require().Equal(sdk.NewInt64DecCoin("uojo", 89), app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddr2).Rewards[0])
}

var historacleTestCases = []struct {
	exchangeRates                         []string
	expectedHistoricMedians               []sdk.Dec
	expectedHistoricMedianDeviation       sdk.Dec
	expectedWithinHistoricMedianDeviation bool
	expectedMedianOfHistoricMedians       sdk.Dec
	expectedAverageOfHistoricMedians      sdk.Dec
	expectedMinOfHistoricMedians          sdk.Dec
	expectedMaxOfHistoricMedians          sdk.Dec
}{
	{
		[]string{
			"1.0", "1.2", "1.1", "1.4", "1.1", "1.15",
			"1.2", "1.3", "1.2", "1.12", "1.2", "1.15",
			"1.17", "1.1", "1.0", "1.16", "1.15", "1.12",
		},
		[]sdk.Dec{
			sdk.MustNewDecFromStr("1.155"),
			sdk.MustNewDecFromStr("1.16"),
			sdk.MustNewDecFromStr("1.175"),
			sdk.MustNewDecFromStr("1.2"),
		},
		sdk.MustNewDecFromStr("0.098615414616580085"),
		true,
		sdk.MustNewDecFromStr("1.1675"),
		sdk.MustNewDecFromStr("1.1725"),
		sdk.MustNewDecFromStr("1.155"),
		sdk.MustNewDecFromStr("1.2"),
	},
	{
		[]string{
			"2.3", "2.12", "2.14", "2.24", "2.18", "2.15",
			"2.51", "2.59", "2.67", "2.76", "2.89", "2.85",
			"3.17", "3.15", "3.35", "3.56", "3.55", "3.49",
		},
		[]sdk.Dec{
			sdk.MustNewDecFromStr("3.02"),
			sdk.MustNewDecFromStr("2.715"),
			sdk.MustNewDecFromStr("2.405"),
			sdk.MustNewDecFromStr("2.24"),
		},
		sdk.MustNewDecFromStr("0.380909000506245145"),
		false,
		sdk.MustNewDecFromStr("2.56"),
		sdk.MustNewDecFromStr("2.595"),
		sdk.MustNewDecFromStr("2.24"),
		sdk.MustNewDecFromStr("3.02"),
	},
	{
		[]string{
			"5.2", "5.25", "5.31", "5.22", "5.14", "5.15",
			"4.85", "4.72", "4.52", "4.47", "4.36", "4.22",
			"4.11", "4.04", "3.92", "3.82", "3.85", "3.83",
		},
		[]sdk.Dec{
			sdk.MustNewDecFromStr("4.165"),
			sdk.MustNewDecFromStr("4.495"),
			sdk.MustNewDecFromStr("4.995"),
			sdk.MustNewDecFromStr("5.15"),
		},
		sdk.MustNewDecFromStr("0.440482689784740573"),
		true,
		sdk.MustNewDecFromStr("4.745"),
		sdk.MustNewDecFromStr("4.70125"),
		sdk.MustNewDecFromStr("4.165"),
		sdk.MustNewDecFromStr("5.15"),
	},
}

func (s *IntegrationTestSuite) TestEndBlockerHistoracle() {
	app, ctx := s.app, s.ctx
	initHeight := ctx.BlockHeight()

	// update historacle params
	app.OracleKeeper.SetHistoricStampPeriod(ctx, 5)
	app.OracleKeeper.SetMedianStampPeriod(ctx, 15)
	app.OracleKeeper.SetMaximumPriceStamps(ctx, 12)
	app.OracleKeeper.SetMaximumMedianStamps(ctx, 4)

	s.T().Log(app.OracleKeeper.MedianStampPeriod(ctx))
	for _, tc := range historacleTestCases {
		ctx = ctx.WithBlockHeight(int64(app.OracleKeeper.MedianStampPeriod(ctx)) - 1)

		for _, exchangeRate := range tc.exchangeRates {
			var tuples types.ExchangeRateTuples
			for _, denom := range app.OracleKeeper.AcceptList(ctx) {
				tuples = append(tuples, types.ExchangeRateTuple{
					Denom:        denom.SymbolDenom,
					ExchangeRate: sdk.MustNewDecFromStr(exchangeRate),
				})
			}

			prevote := types.AggregateExchangeRatePrevote{
				Hash:        "hash",
				Voter:       valAddr1.String(),
				SubmitBlock: uint64(ctx.BlockHeight()),
			}
			app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr1, prevote)
			oracle.EndBlocker(ctx, app.OracleKeeper)

			ctx = ctx.WithBlockHeight(ctx.BlockHeight() + int64(app.OracleKeeper.VotePeriod(ctx)))
			vote := types.AggregateExchangeRateVote{
				ExchangeRateTuples: tuples,
				Voter:              valAddr1.String(),
			}
			app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr1, vote)
			oracle.EndBlocker(ctx, app.OracleKeeper)
		}

		for _, denom := range app.OracleKeeper.AcceptList(ctx) {
			// query for past 6 medians (should only get 4 back since max median stamps is set to 4)
			medians := app.OracleKeeper.HistoricMedians(ctx, denom.SymbolDenom, 6)
			s.Require().Equal(4, len(medians))
			s.Require().Equal(tc.expectedHistoricMedians, medians)

			medianHistoricDeviation, err := app.OracleKeeper.HistoricMedianDeviation(ctx, denom.SymbolDenom)
			s.Require().NoError(err)
			s.Require().Equal(tc.expectedHistoricMedianDeviation, medianHistoricDeviation)

			withinHistoricMedianDeviation, err := app.OracleKeeper.WithinHistoricMedianDeviation(ctx, denom.SymbolDenom)
			s.Require().NoError(err)
			s.Require().Equal(tc.expectedWithinHistoricMedianDeviation, withinHistoricMedianDeviation)

			medianOfHistoricMedians, numMedians, err := app.OracleKeeper.MedianOfHistoricMedians(ctx, denom.SymbolDenom, 6)
			s.Require().Equal(uint32(4), numMedians)
			s.Require().Equal(tc.expectedMedianOfHistoricMedians, medianOfHistoricMedians)

			averageOfHistoricMedians, numMedians, err := app.OracleKeeper.AverageOfHistoricMedians(ctx, denom.SymbolDenom, 6)
			s.Require().Equal(uint32(4), numMedians)
			s.Require().Equal(tc.expectedAverageOfHistoricMedians, averageOfHistoricMedians)

			minOfHistoricMedians, numMedians, err := app.OracleKeeper.MinOfHistoricMedians(ctx, denom.SymbolDenom, 6)
			s.Require().Equal(uint32(4), numMedians)
			s.Require().Equal(tc.expectedMinOfHistoricMedians, minOfHistoricMedians)

			maxOfHistoricMedians, numMedians, err := app.OracleKeeper.MaxOfHistoricMedians(ctx, denom.SymbolDenom, 6)
			s.Require().Equal(uint32(4), numMedians)
			s.Require().Equal(tc.expectedMaxOfHistoricMedians, maxOfHistoricMedians)

			clearHistoricPrices(ctx, app.OracleKeeper, denom.SymbolDenom)
			clearHistoricMedians(ctx, app.OracleKeeper, denom.SymbolDenom)
			clearHistoricMedianDeviations(ctx, app.OracleKeeper, denom.SymbolDenom)
		}

		ctx = ctx.WithBlockHeight(initHeight)
	}
}

func TestOracleTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
