package abci_test

import (
	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"

	"github.com/ojo-network/ojo/pricefeeder"
	"github.com/ojo-network/ojo/x/oracle/abci"
	"github.com/ojo-network/ojo/x/oracle/keeper"
)

func (s *IntegrationTestSuite) TestExtendVoteHandler() {
	app, ctx := s.app, s.ctx

	testCases := []struct {
		name              string
		logger            log.Logger
		oracleKeeper      keeper.Keeper
		priceFeeder       *pricefeeder.PriceFeeder
		extendVoteRequest func() *cometabci.RequestExtendVote
		expErr            bool
		expErrMsg         string
	}{
		{
			name:         "nil price feeder",
			logger:       log.NewNopLogger(),
			oracleKeeper: app.OracleKeeper,
			priceFeeder:  app.PriceFeeder,
			expErr:       true,
			expErrMsg:    "price feeder oracle not set",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			h := abci.NewVoteExtensionHandler(
				tc.logger,
				tc.oracleKeeper,
				tc.priceFeeder,
			)

			req := &cometabci.RequestExtendVote{}
			if tc.extendVoteRequest != nil {
				req = tc.extendVoteRequest()
			}
			resp, err := h.ExtendVoteHandler()(ctx, req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				if resp == nil || len(resp.VoteExtension) == 0 {
					return
				}
				s.Require().NoError(err)
				s.Require().NotNil(resp)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestVerifyVoteExtensionHandler() {

}
