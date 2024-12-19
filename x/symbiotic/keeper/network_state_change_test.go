package keeper_test

import (
	"time"

	"github.com/ojo-network/ojo/x/symbiotic/types"
)

const (
	ojoMiddlewareAddress = "0x127E303D6604C48f3bA0010EbEa57e09324A4dF6"
)

func (s *IntegrationTestSuite) TestSymbioticUpdateValidatorsPower() {
	app, ctx := s.app, s.ctx
	ctx = ctx.
		WithBlockHeight(10).
		WithBlockTime(time.Now())

	params := app.SymbioticKeeper.GetParams(ctx)
	params.MiddlewareAddress = ojoMiddlewareAddress
	app.SymbioticKeeper.SetParams(ctx, params)

	blockHash, err := app.SymbioticKeeper.GetFinalizedBlockHash(ctx)
	s.Require().NoError(err)

	cachedBlockHash := types.CachedBlockHash{
		BlockHash: blockHash,
		Height:    ctx.BlockHeight(),
	}

	s.app.SymbioticKeeper.SetCachedBlockHash(ctx, cachedBlockHash)

	err = s.app.SymbioticKeeper.SymbioticUpdateValidatorsPower(ctx)
	s.Require().NoError(err)
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
