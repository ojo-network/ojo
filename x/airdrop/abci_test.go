package airdrop_test

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
	"github.com/ojo-network/ojo/x/airdrop"
	"github.com/ojo-network/ojo/x/airdrop/types"
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
	initialPower = int64(1000)
)

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupTest() {
	require := s.Require()
	isCheckTx := false
	app := ojoapp.Setup(s.T())
	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{
		ChainID: fmt.Sprintf("test-chain-%s", tmrand.Str(4)),
	})

	airdrop.InitGenesis(ctx, app.AirdropKeeper, *types.DefaultGenesisState())

	setupVals := app.StakingKeeper.GetBondedValidatorsByPower(ctx)
	s.Require().Len(setupVals, 1)
	s.Require().Equal(int64(1), setupVals[0].GetConsensusPower(app.StakingKeeper.PowerReduction(ctx)))

	sh := teststaking.NewHelper(s.T(), ctx, *app.StakingKeeper)
	sh.Denom = bondDenom

	// mint and send coins to validators
	require.NoError(app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initCoins.MulInt(sdk.NewIntFromUint64(3))))
	require.NoError(app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr1, initCoins))
	require.NoError(app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr2, initCoins))
	require.NoError(app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr3, initCoins))

	// mint and send coins to oracle module to fill up reward pool
	require.NoError(app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initCoins))
	require.NoError(app.BankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, types.ModuleName, initCoins))

	sh.CreateValidatorWithValPower(valAddr1, valPubKey1, 599, true)
	sh.CreateValidatorWithValPower(valAddr2, valPubKey2, 398, true)
	sh.CreateValidatorWithValPower(valAddr3, valPubKey3, 2, true)

	staking.EndBlocker(ctx, *app.StakingKeeper)

	s.app = app
	s.ctx = ctx
}

// Test addresses
var (
	valPubKeys = simapp.CreateTestPubKeys(3)

	valPubKey1 = valPubKeys[0]
	pubKey1    = secp256k1.GenPrivKey().PubKey()
	addr1      = sdk.AccAddress(pubKey1.Address())
	valAddr1   = sdk.ValAddress(pubKey1.Address())

	valPubKey2 = valPubKeys[1]
	pubKey2    = secp256k1.GenPrivKey().PubKey()
	addr2      = sdk.AccAddress(pubKey2.Address())
	valAddr2   = sdk.ValAddress(pubKey2.Address())

	valPubKey3 = valPubKeys[2]
	pubKey3    = secp256k1.GenPrivKey().PubKey()
	addr3      = sdk.AccAddress(pubKey3.Address())
	valAddr3   = sdk.ValAddress(pubKey3.Address())

	initTokens = sdk.TokensFromConsensusPower(initialPower, sdk.DefaultPowerReduction)
	initCoins  = sdk.NewCoins(sdk.NewCoin(bondDenom, initTokens))
)

func (s *IntegrationTestSuite) TestEndBlockerAccountCreation() {
	app, ctx := s.app, s.ctx

	pubKey := secp256k1.GenPrivKey().PubKey()
	testAddress := sdk.AccAddress(pubKey.Address())

	originAmount := uint64(600)

	airdropAccount := &types.AirdropAccount{
		OriginAddress:  testAddress.String(),
		OriginAmount:   originAmount,
		VestingEndTime: int64(10000000000),
		State:          types.AirdropAccount_CREATED,
	}
	app.AirdropKeeper.SetAirdropAccount(ctx, airdropAccount)

	airdrop.EndBlocker(ctx, app.AirdropKeeper)

	accAddress, err := airdropAccount.OriginAccAddress()
	s.Require().NoError(err)

	balance := app.BankKeeper.GetBalance(
		ctx,
		accAddress,
		bondDenom,
	).Amount

	s.Require().Equal(originAmount, balance.Uint64())

}

func (s *IntegrationTestSuite) TestEndBlockerMinting() {
	app, ctx := s.app, s.ctx

	distributionStartingBalance := app.BankKeeper.GetBalance(
		ctx,
		app.AirdropKeeper.DistributionModuleAddress(ctx),
		bondDenom,
	).Amount

	communityPoolStartingBalance := s.app.DistrKeeper.GetFeePool(ctx).CommunityPool.AmountOf(bondDenom)

	airdropAccount := &types.AirdropAccount{
		OriginAddress:  "testAddress",
		OriginAmount:   uint64(600),
		VestingEndTime: int64(10000000000),
		State:          types.AirdropAccount_UNCLAIMED,
	}
	app.AirdropKeeper.SetAirdropAccount(ctx, airdropAccount)

	ctx = ctx.WithBlockHeight(int64(app.AirdropKeeper.GetParams(ctx).ExpiryBlock))
	airdrop.EndBlocker(ctx, app.AirdropKeeper)

	// Check that the airdrop account has been claimed
	queriedAccount, err := app.AirdropKeeper.GetAirdropAccount(ctx, airdropAccount.OriginAddress, types.AirdropAccount_CLAIMED)
	s.Require().NoError(err)
	err = queriedAccount.VerifyNotClaimed()
	s.Require().Error(err)

	// Check that the distribution module account has increased by the claim amount
	distributionEndingBalance := app.BankKeeper.GetBalance(
		ctx,
		app.AirdropKeeper.DistributionModuleAddress(ctx),
		bondDenom,
	).Amount
	distributionDifference := distributionEndingBalance.Sub(distributionStartingBalance)
	s.Require().Equal(queriedAccount.ClaimAmount, distributionDifference.Uint64())

	// Check that the community pool balance has been updated
	communityPoolEndingBalance := s.app.DistrKeeper.GetFeePool(ctx).CommunityPool.AmountOf(bondDenom)
	communityPoolDifference := communityPoolEndingBalance.Sub(communityPoolStartingBalance)
	s.Require().Equal(queriedAccount.ClaimAmount, uint64(communityPoolDifference.TruncateInt64()))
}
