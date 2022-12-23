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
	displayDenom string = appparams.DisplayDenom
	bondDenom    string = appparams.BondDenom
)

type IntegrationTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *ojoapp.App
}

const (
	initialPower = int64(10000000000)
)

func (s *IntegrationTestSuite) SetupTest() {
	require := s.Require()
	isCheckTx := false
	app := ojoapp.Setup(s.T(), false, 1)
	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{
		ChainID: fmt.Sprintf("test-chain-%s", tmrand.Str(4)),
		Height:  9,
	})

	oracle.InitGenesis(ctx, app.OracleKeeper, *types.DefaultGenesisState())

	sh := teststaking.NewHelper(s.T(), ctx, *app.StakingKeeper)
	sh.Denom = bondDenom
	amt1 := sdk.TokensFromConsensusPower(70, sdk.DefaultPowerReduction)
	amt2 := sdk.TokensFromConsensusPower(30, sdk.DefaultPowerReduction)

	// mint and send coins to validator
	require.NoError(app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initCoins))
	require.NoError(app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr1, initCoins))
	require.NoError(app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initCoins))
	require.NoError(app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr2, initCoins))

	sh.CreateValidator(valAddr1, valPubKey1, amt1, true)
	sh.CreateValidator(valAddr2, valPubKey2, amt2, true)

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

func (s *IntegrationTestSuite) TestEnblockerVoteThreshold() {
	app, ctx := s.app, s.ctx

	var val1Tuples types.ExchangeRateTuples
	var val2Tuples types.ExchangeRateTuples
	var val1PreVotes []types.AggregateExchangeRatePrevote
	var val2PreVotes []types.AggregateExchangeRatePrevote
	var val1Votes []types.AggregateExchangeRateVote
	var val2Votes []types.AggregateExchangeRateVote
	for _, denom := range app.OracleKeeper.AcceptList(ctx) {
		val1Tuples = append(val1Tuples, types.ExchangeRateTuple{
			Denom:        denom.SymbolDenom,
			ExchangeRate: sdk.MustNewDecFromStr("1.0"),
		})
		val2Tuples = append(val2Tuples, types.ExchangeRateTuple{
			Denom:        denom.SymbolDenom,
			ExchangeRate: sdk.MustNewDecFromStr("0.5"),
		})

		val1PreVotes = append(val1PreVotes, types.AggregateExchangeRatePrevote{
			Hash:        "hash1",
			Voter:       valAddr1.String(),
			SubmitBlock: uint64(ctx.BlockHeight()),
		})
		val2PreVotes = append(val2PreVotes, types.AggregateExchangeRatePrevote{
			Hash:        "hash2",
			Voter:       valAddr2.String(),
			SubmitBlock: uint64(ctx.BlockHeight()),
		})

		val1Votes = append(val1Votes, types.AggregateExchangeRateVote{
			ExchangeRateTuples: val1Tuples,
			Voter:              valAddr1.String(),
		})
		val2Votes = append(val2Votes, types.AggregateExchangeRateVote{
			ExchangeRateTuples: val2Tuples,
			Voter:              valAddr2.String(),
		})
	}

	// total voting power per denom is 100
	for i := range val1PreVotes {
		app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr1, val1PreVotes[i])
		app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr2, val2PreVotes[i])
	}
	oracle.EndBlocker(ctx, app.OracleKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + int64(app.OracleKeeper.VotePeriod(ctx)))
	for i := range val1Votes {
		app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr1, val1Votes[i])
		app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr2, val2Votes[i])
	}
	oracle.EndBlocker(ctx, app.OracleKeeper)

	for _, denom := range app.OracleKeeper.AcceptList(ctx) {
		rate, err := app.OracleKeeper.GetExchangeRate(ctx, denom.SymbolDenom)
		s.Require().NoError(err)
		s.Require().Equal(sdk.MustNewDecFromStr("0.75"), rate)
	}

	// total voting power per denom is 30
	for i := range val2PreVotes {
		app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr2, val2PreVotes[i])
	}
	oracle.EndBlocker(ctx, app.OracleKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + int64(app.OracleKeeper.VotePeriod(ctx)))
	for i := range val2Votes {
		app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr2, val2Votes[i])
	}
	oracle.EndBlocker(ctx, app.OracleKeeper)

	for _, denom := range app.OracleKeeper.AcceptList(ctx) {
		rate, err := app.OracleKeeper.GetExchangeRate(ctx, denom.SymbolDenom)
		s.Require().ErrorIs(err, sdkerrors.Wrap(types.ErrUnknownDenom, denom.SymbolDenom))
		s.Require().Equal(sdk.ZeroDec(), rate)
	}
}

func TestOracleTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
