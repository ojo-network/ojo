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
	"github.com/ojo-network/ojo/x/oracle/types"
)

const (
	displayDenom     string = appparams.DisplayDenom
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

func (s *IntegrationTestSuite) SetupTest() {
	require := s.Require()
	isCheckTx := false
	app := ojoapp.Setup(s.T(), false, 1)
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

func (s *IntegrationTestSuite) TestEndblockerVoteThreshold() {
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

func (s *IntegrationTestSuite) TestEndblockerValidatorRewards() {
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

func TestOracleTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
