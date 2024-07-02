package abci_test

import (
	"encoding/json"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	cometabci "github.com/cometbft/cometbft/abci/types"
	cometprototypes "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/ojo-network/ojo/x/oracle/abci"
	"github.com/ojo-network/ojo/x/oracle/keeper"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

func (s *IntegrationTestSuite) TestPreBlocker() {
	app, ctx, keys := s.app, s.ctx, s.keys
	voter := keys[0].ValAddress

	// enable vote extensions
	ctx = ctx.WithBlockHeight(3)
	consensusParams := ctx.ConsensusParams()
	consensusParams.Abci = &cometprototypes.ABCIParams{
		VoteExtensionsEnableHeight: 2,
	}
	ctx = ctx.WithConsensusParams(consensusParams)

	// build injected vote extention tx
	exchangeRateVoteAtom := oracletypes.AggregateExchangeRateVote{
		ExchangeRates: sdk.NewDecCoinsFromCoins(sdk.NewCoin("ATOM", math.NewInt(11))),
		Voter:         voter.String(),
	}
	injectedVoteExtTx := abci.AggregateExchangeRateVotes{
		ExchangeRateVotes: []oracletypes.AggregateExchangeRateVote{
			exchangeRateVoteAtom,
		},
	}
	bz, err := json.Marshal(injectedVoteExtTx)
	s.Require().NoError(err)
	var txs [][]byte
	txs = append(txs, bz)

	testCases := []struct {
		name                 string
		logger               log.Logger
		oracleKeeper         keeper.Keeper
		finalizeBlockRequest *cometabci.RequestFinalizeBlock
		expErr               bool
		expErrMsg            string
	}{
		{
			name:         "nil preblocker request",
			logger:       log.NewNopLogger(),
			oracleKeeper: app.OracleKeeper,
			expErr:       true,
			expErrMsg:    "preblocker received a nil request",
		},
		{
			name:         "oracle preblocker sets exchange rate",
			logger:       log.NewNopLogger(),
			oracleKeeper: app.OracleKeeper,
			finalizeBlockRequest: &cometabci.RequestFinalizeBlock{
				Height: ctx.BlockHeight(),
				Txs:    txs,
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			h := abci.NewPreBlockHandler(
				tc.logger,
				tc.oracleKeeper,
			)

			resp, err := h.PreBlocker(module.NewManager())(ctx, tc.finalizeBlockRequest)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)

				// check exchange rate vote was set
				exchangeRateVote, err := tc.oracleKeeper.GetAggregateExchangeRateVote(ctx, voter)
				s.Require().NoError(err)
				s.Require().Equal(exchangeRateVoteAtom, exchangeRateVote)
			}
		})
	}
}
