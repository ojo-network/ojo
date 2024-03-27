package abci_test

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"

	"github.com/ojo-network/ojo/pricefeeder"
	"github.com/ojo-network/ojo/x/oracle/abci"
	"github.com/ojo-network/ojo/x/oracle/keeper"
)

func (s *IntegrationTestSuite) TestExtendVoteHandler() {
	app, ctx := s.app, s.ctx
	pf := MockPriceFeeder()

	testCases := []struct {
		name              string
		logger            log.Logger
		oracleKeeper      keeper.Keeper
		priceFeeder       *pricefeeder.PriceFeeder
		extendVoteRequest *cometabci.RequestExtendVote
		expErr            bool
		expErrMsg         string
	}{
		{
			name:         "nil vote extension request",
			logger:       log.NewNopLogger(),
			oracleKeeper: app.OracleKeeper,
			priceFeeder:  pf,
			expErr:       true,
			expErrMsg:    "extend vote handler received a nil request",
		},
		{
			name:         "price feeder oracle not set",
			logger:       log.NewNopLogger(),
			oracleKeeper: app.OracleKeeper,
			priceFeeder:  app.PriceFeeder,
			extendVoteRequest: &cometabci.RequestExtendVote{
				Height: ctx.BlockHeight(),
			},
			expErr:    true,
			expErrMsg: "price feeder oracle not set",
		},
		{
			name:         "vote extension handled successfully",
			logger:       log.NewNopLogger(),
			oracleKeeper: app.OracleKeeper,
			priceFeeder:  pf,
			extendVoteRequest: &cometabci.RequestExtendVote{
				Height: ctx.BlockHeight(),
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			h := abci.NewVoteExtensionHandler(
				tc.logger,
				tc.oracleKeeper,
				tc.priceFeeder,
			)

			resp, err := h.ExtendVoteHandler()(ctx, tc.extendVoteRequest)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				s.Require().Greater(len(resp.VoteExtension), 0)

				var voteExt abci.OracleVoteExtension
				err = json.Unmarshal(resp.VoteExtension, &voteExt)
				s.Require().NoError(err)
				s.Require().Equal(ctx.BlockHeight(), voteExt.Height)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestVerifyVoteExtensionHandler() {
	app, ctx := s.app, s.ctx
	pf := MockPriceFeeder()

	voteExtension, err := json.Marshal(&cometabci.RequestExtendVote{
		Height: ctx.BlockHeight(),
	})
	s.Require().NoError(err)

	testCases := []struct {
		name              string
		logger            log.Logger
		oracleKeeper      keeper.Keeper
		priceFeeder       *pricefeeder.PriceFeeder
		verifyVoteRequest *cometabci.RequestVerifyVoteExtension
		expErr            bool
		expErrMsg         string
	}{
		{
			name:         "nil verify vote extension request",
			logger:       log.NewNopLogger(),
			oracleKeeper: app.OracleKeeper,
			priceFeeder:  pf,
			expErr:       true,
			expErrMsg:    "verify vote extension handler received a nil request",
		},
		{
			name:         "vote extension and verify vote extension height mismatch",
			logger:       log.NewNopLogger(),
			oracleKeeper: app.OracleKeeper,
			priceFeeder:  pf,
			verifyVoteRequest: &cometabci.RequestVerifyVoteExtension{
				Height:        ctx.BlockHeight() + 1,
				VoteExtension: voteExtension,
			},
			expErr: true,
			expErrMsg: fmt.Sprintf("verify vote extension handler received vote extension height that doesn't"+
				"match request height; expected: %d, got: %d",
				ctx.BlockHeight()+1,
				ctx.BlockHeight(),
			),
		},
		{
			name:         "vote extension verified successfully",
			logger:       log.NewNopLogger(),
			oracleKeeper: app.OracleKeeper,
			priceFeeder:  pf,
			verifyVoteRequest: &cometabci.RequestVerifyVoteExtension{
				Height:        ctx.BlockHeight(),
				VoteExtension: voteExtension,
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			h := abci.NewVoteExtensionHandler(
				tc.logger,
				tc.oracleKeeper,
				tc.priceFeeder,
			)

			resp, err := h.VerifyVoteExtensionHandler()(ctx, tc.verifyVoteRequest)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				s.Require().Equal(cometabci.ResponseVerifyVoteExtension_ACCEPT, resp.Status)
			}
		})
	}
}
