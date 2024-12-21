package keeper_test

import (
	"time"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtestutil "github.com/cosmos/cosmos-sdk/x/staking/testutil"
	appparams "github.com/ojo-network/ojo/app/params"
	"github.com/ojo-network/ojo/x/symbiotic/types"
)

const (
	ojoMiddlewareAddress = "0x127E303D6604C48f3bA0010EbEa57e09324A4dF6"
	secret               = "fixed_seed_value"
)

func (s *IntegrationTestSuite) TestSymbioticUpdateValidatorsPower() {
	app, ctx := s.app, s.ctx
	ctx = ctx.
		WithBlockHeight(10).
		WithBlockTime(time.Now())

	sh := stakingtestutil.NewHelper(s.T(), ctx, app.StakingKeeper)
	sh.Denom = appparams.BondDenom
	privKey := secp256k1.GenPrivKeyFromSecret([]byte(secret))
	pubKey := privKey.PubKey()

	initTokens := sdk.TokensFromConsensusPower(1000, sdk.DefaultPowerReduction)
	initCoins := sdk.NewCoins(sdk.NewCoin(appparams.BondDenom, initTokens))

	s.app.BankKeeper.MintCoins(s.ctx, minttypes.ModuleName, initCoins)
	s.app.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, minttypes.ModuleName, sdk.AccAddress(pubKey.Address()), initCoins)

	sh.CreateValidatorWithValPower(
		sdk.ValAddress(pubKey.Address()),
		pubKey,
		1000,
		true,
	)
	val, err := app.StakingKeeper.GetValidator(ctx, sdk.ValAddress(pubKey.Address()))
	s.Require().Equal(val.GetTokens(), math.NewInt(1000000000))

	params := app.SymbioticKeeper.GetParams(ctx)
	params.MiddlewareAddress = ojoMiddlewareAddress
	app.SymbioticKeeper.SetParams(ctx, params)

	blockHash, err := app.SymbioticKeeper.GetFinalizedBlockHash(ctx)
	s.Require().NoError(err)

	cachedBlockHash := types.CachedBlockHash{
		BlockHash: blockHash,
		Height:    ctx.BlockHeight(),
	}

	app.SymbioticKeeper.SetCachedBlockHash(ctx, cachedBlockHash)

	err = app.SymbioticKeeper.SymbioticUpdateValidatorsPower(ctx)
	s.Require().NoError(err)

	val, err = app.StakingKeeper.GetValidator(ctx, sdk.ValAddress(pubKey.Address()))
	s.Require().NoError(err)
	s.Require().Equal(val.GetTokens(), math.NewInt(0))
}

func (s *IntegrationTestSuite) TestFinalizedBlockHash() {
	app, ctx := s.app, s.ctx
	ctx = ctx.
		WithBlockHeight(10).
		WithBlockTime(time.Now())

	blockHash, err := app.SymbioticKeeper.GetFinalizedBlockHash(ctx)
	s.Require().NoError(err)

	_, err = app.SymbioticKeeper.GetBlockByHash(ctx, blockHash)
	s.Require().NoError(err)
}
