package keeper_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cometbft/cometbft/crypto/secp256k1"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtestutil "github.com/cosmos/cosmos-sdk/x/staking/testutil"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"

	ojoapp "github.com/ojo-network/ojo/app"
	appparams "github.com/ojo-network/ojo/app/params"
	"github.com/ojo-network/ojo/x/oracle/keeper"
	"github.com/ojo-network/ojo/x/oracle/types"
)

const (
	displayDenom string = appparams.DisplayDenom
	bondDenom    string = appparams.BondDenom
)

type IntegrationTestSuite struct {
	suite.Suite

	ctx         sdk.Context
	app         *ojoapp.App
	queryClient types.QueryClient
	msgServer   types.MsgServer
}

const (
	initialPower = int64(10000000000)
)

func (s *IntegrationTestSuite) SetupTest() {
	// set prefixes for valid addresses, without sealing
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(appparams.AccountAddressPrefix, appparams.AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(appparams.ValidatorAddressPrefix, appparams.ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(appparams.ConsNodeAddressPrefix, appparams.ConsNodePubKeyPrefix)

	require := s.Require()
	isCheckTx := false
	app := ojoapp.Setup(s.T())
	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{
		ChainID: fmt.Sprintf("test-chain-%s", tmrand.Str(4)),
		Height:  9,
	})

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, keeper.NewQuerier(app.OracleKeeper))

	sh := stakingtestutil.NewHelper(s.T(), ctx, app.StakingKeeper)
	sh.Denom = bondDenom
	amt := sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)

	// mint and send coins to validators
	require.NoError(app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initCoins))
	require.NoError(app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, initCoins))
	require.NoError(app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initCoins))
	require.NoError(app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr2, initCoins))

	sh.CreateValidator(valAddr, valPubKey, amt, true)
	sh.CreateValidator(valAddr2, valPubKey2, amt, true)

	staking.EndBlocker(ctx, app.StakingKeeper)

	s.app = app
	s.ctx = ctx
	s.queryClient = types.NewQueryClient(queryHelper)
	s.msgServer = keeper.NewMsgServerImpl(app.OracleKeeper)
}

// Test addresses
var (
	valPubKeys = simtestutil.CreateTestPubKeys(2)

	valPubKey = valPubKeys[0]
	pubKey    = secp256k1.GenPrivKey().PubKey()
	addr      = sdk.AccAddress(pubKey.Address())
	valAddr   = sdk.ValAddress(pubKey.Address())

	valPubKey2 = valPubKeys[1]
	pubKey2    = secp256k1.GenPrivKey().PubKey()
	addr2      = sdk.AccAddress(pubKey2.Address())
	valAddr2   = sdk.ValAddress(pubKey2.Address())

	initTokens = sdk.TokensFromConsensusPower(initialPower, sdk.DefaultPowerReduction)
	initCoins  = sdk.NewCoins(sdk.NewCoin(appparams.BondDenom, initTokens))
)

// NewTestMsgCreateValidator test msg creator
func NewTestMsgCreateValidator(address sdk.ValAddress, pubKey cryptotypes.PubKey, amt sdk.Int) *stakingtypes.MsgCreateValidator {
	commission := stakingtypes.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
	msg, _ := stakingtypes.NewMsgCreateValidator(
		address, pubKey, sdk.NewCoin(types.OjoDenom, amt),
		stakingtypes.Description{}, commission, sdk.OneInt(),
	)

	return msg
}

func (s *IntegrationTestSuite) TestSetFeederDelegation() {
	app, ctx := s.app, s.ctx

	feederAddr := sdk.AccAddress([]byte("addr________________"))
	feederAcc := app.AccountKeeper.NewAccountWithAddress(ctx, feederAddr)
	app.AccountKeeper.SetAccount(ctx, feederAcc)

	err := s.app.OracleKeeper.ValidateFeeder(ctx, addr, valAddr)
	s.Require().NoError(err)
	err = s.app.OracleKeeper.ValidateFeeder(ctx, feederAddr, valAddr)
	s.Require().Error(err)

	s.app.OracleKeeper.SetFeederDelegation(ctx, valAddr, feederAddr)

	err = s.app.OracleKeeper.ValidateFeeder(ctx, addr, valAddr)
	s.Require().Error(err)
	err = s.app.OracleKeeper.ValidateFeeder(ctx, feederAddr, valAddr)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TestGetFeederDelegation() {
	app, ctx := s.app, s.ctx

	feederAddr := sdk.AccAddress([]byte("addr________________"))
	feederAcc := app.AccountKeeper.NewAccountWithAddress(ctx, feederAddr)
	app.AccountKeeper.SetAccount(ctx, feederAcc)

	s.app.OracleKeeper.SetFeederDelegation(ctx, valAddr, feederAddr)
	resp, err := app.OracleKeeper.GetFeederDelegation(ctx, valAddr)
	s.Require().NoError(err)
	s.Require().Equal(resp, feederAddr)
}

func (s *IntegrationTestSuite) TestAggregateExchangeRatePrevote() {
	app, ctx := s.app, s.ctx

	prevote := types.AggregateExchangeRatePrevote{
		Hash:        "hash",
		Voter:       addr.String(),
		SubmitBlock: 0,
	}
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr, prevote)

	_, err := app.OracleKeeper.GetAggregateExchangeRatePrevote(ctx, valAddr)
	s.Require().NoError(err)

	app.OracleKeeper.DeleteAggregateExchangeRatePrevote(ctx, valAddr)

	_, err = app.OracleKeeper.GetAggregateExchangeRatePrevote(ctx, valAddr)
	s.Require().Error(err)
}

func (s *IntegrationTestSuite) TestAggregateExchangeRatePrevoteError() {
	app, ctx := s.app, s.ctx

	_, err := app.OracleKeeper.GetAggregateExchangeRatePrevote(ctx, valAddr)
	s.Require().Errorf(err, types.ErrNoAggregatePrevote.Error())
}

func (s *IntegrationTestSuite) TestAggregateExchangeRateVote() {
	app, ctx := s.app, s.ctx

	var decCoins sdk.DecCoins
	decCoins = append(decCoins, sdk.DecCoin{
		Denom:  displayDenom,
		Amount: sdk.ZeroDec(),
	})

	vote := types.AggregateExchangeRateVote{
		ExchangeRates: decCoins,
		Voter:         addr.String(),
	}
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr, vote)

	_, err := app.OracleKeeper.GetAggregateExchangeRateVote(ctx, valAddr)
	s.Require().NoError(err)

	app.OracleKeeper.DeleteAggregateExchangeRateVote(ctx, valAddr)

	_, err = app.OracleKeeper.GetAggregateExchangeRateVote(ctx, valAddr)
	s.Require().Error(err)
}

func (s *IntegrationTestSuite) TestAggregateExchangeRateVoteError() {
	app, ctx := s.app, s.ctx

	_, err := app.OracleKeeper.GetAggregateExchangeRateVote(ctx, valAddr)
	s.Require().Errorf(err, types.ErrNoAggregateVote.Error())
}

func (s *IntegrationTestSuite) TestSetExchangeRateWithEvent() {
	app, ctx := s.app, s.ctx
	err := app.OracleKeeper.SetExchangeRateWithEvent(ctx, displayDenom, sdk.OneDec())
	s.Require().NoError(err)
	rate, err := app.OracleKeeper.GetExchangeRate(ctx, displayDenom)
	s.Require().NoError(err)
	s.Require().Equal(rate, sdk.OneDec())
}

func (s *IntegrationTestSuite) TestGetExchangeRate_InvalidDenom() {
	app, ctx := s.app, s.ctx

	_, err := app.OracleKeeper.GetExchangeRate(ctx, "uxyz")
	s.Require().Error(err)
}

func (s *IntegrationTestSuite) TestGetExchangeRate_NotSet() {
	app, ctx := s.app, s.ctx

	_, err := app.OracleKeeper.GetExchangeRate(ctx, displayDenom)
	s.Require().Error(err)
}

func (s *IntegrationTestSuite) TestGetExchangeRate_Valid() {
	app, ctx := s.app, s.ctx

	app.OracleKeeper.SetExchangeRate(ctx, displayDenom, sdk.OneDec())
	rate, err := app.OracleKeeper.GetExchangeRate(ctx, displayDenom)
	s.Require().NoError(err)
	s.Require().Equal(rate, sdk.OneDec())

	app.OracleKeeper.SetExchangeRate(ctx, strings.ToLower(displayDenom), sdk.OneDec())
	rate, err = app.OracleKeeper.GetExchangeRate(ctx, displayDenom)
	s.Require().NoError(err)
	s.Require().Equal(rate, sdk.OneDec())
}

func (s *IntegrationTestSuite) TestGetExchangeRateBase() {
	oracleParams := s.app.OracleKeeper.GetParams(s.ctx)

	var exponent uint64
	for _, denom := range oracleParams.AcceptList {
		if denom.BaseDenom == bondDenom {
			exponent = uint64(denom.Exponent)
		}
	}

	power := sdk.MustNewDecFromStr("10").Power(exponent)

	s.app.OracleKeeper.SetExchangeRate(s.ctx, displayDenom, sdk.OneDec())
	rate, err := s.app.OracleKeeper.GetExchangeRateBase(s.ctx, bondDenom)
	s.Require().NoError(err)
	s.Require().Equal(rate.Mul(power), sdk.OneDec())

	s.app.OracleKeeper.SetExchangeRate(s.ctx, strings.ToLower(displayDenom), sdk.OneDec())
	rate, err = s.app.OracleKeeper.GetExchangeRateBase(s.ctx, bondDenom)
	s.Require().NoError(err)
	s.Require().Equal(rate.Mul(power), sdk.OneDec())
}

func (s *IntegrationTestSuite) TestClearExchangeRate() {
	app, ctx := s.app, s.ctx

	app.OracleKeeper.SetExchangeRate(ctx, displayDenom, sdk.OneDec())
	app.OracleKeeper.ClearExchangeRates(ctx)
	_, err := app.OracleKeeper.GetExchangeRate(ctx, displayDenom)
	s.Require().Error(err)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
