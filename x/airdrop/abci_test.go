package airdrop_test

import (
	"testing"

	"github.com/cometbft/cometbft/crypto/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	ojoapp "github.com/ojo-network/ojo/app"
	appparams "github.com/ojo-network/ojo/app/params"
	"github.com/ojo-network/ojo/tests/integration"
	"github.com/ojo-network/ojo/x/airdrop"
	"github.com/ojo-network/ojo/x/airdrop/types"
)

const bondDenom = appparams.BondDenom

type IntegrationTestSuite struct {
	suite.Suite

	ctx  sdk.Context
	app  *ojoapp.App
	keys []integration.TestValidatorKey
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupTest() {
	s.app, s.ctx, s.keys = integration.SetupAppWithContext(s.T(), 2)
}

func (s *IntegrationTestSuite) TestEndBlockerAccountCreation() {
	app, ctx := s.app, s.ctx

	pubKey := secp256k1.GenPrivKey().PubKey()
	testAddress := sdk.AccAddress(pubKey.Address())

	originAmount := uint64(600)

	airdropAccount := &types.AirdropAccount{
		OriginAddress:  testAddress.String(),
		OriginAmount:   originAmount,
		VestingEndTime: int64(10000000000),
		State:          types.AirdropAccount_STATE_CREATED,
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
		State:          types.AirdropAccount_STATE_UNCLAIMED,
	}
	app.AirdropKeeper.SetAirdropAccount(ctx, airdropAccount)

	ctx = ctx.WithBlockHeight(int64(app.AirdropKeeper.GetParams(ctx).ExpiryBlock))
	airdrop.EndBlocker(ctx, app.AirdropKeeper)

	// Check that the airdrop account has been claimed
	queriedAccount, err := app.AirdropKeeper.GetAirdropAccount(ctx, airdropAccount.OriginAddress, types.AirdropAccount_STATE_CLAIMED)
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
