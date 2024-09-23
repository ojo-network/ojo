package abci_test

import (
	"bytes"
	"sort"

	"cosmossdk.io/core/header"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	cometabci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protoio "github.com/cosmos/gogoproto/io"
	"github.com/cosmos/gogoproto/proto"
	"github.com/ojo-network/ojo/tests/integration"

	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/ojo-network/ojo/x/oracle/abci"
	oraclekeeper "github.com/ojo-network/ojo/x/oracle/keeper"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

func (s *IntegrationTestSuite) TestPrepareProposalHandler() {
	app, ctx, keys := s.app, s.ctx, s.keys

	// enable vote extensions
	ctx = ctx.WithBlockHeight(3)
	consensusParams := ctx.ConsensusParams()
	consensusParams.Abci = &cmtproto.ABCIParams{
		VoteExtensionsEnableHeight: 2,
	}
	ctx = ctx.WithConsensusParams(consensusParams)

	// build local commit info
	exchangeRates := sdk.NewDecCoinsFromCoins(sdk.NewCoin("ATOM", math.NewInt(11)))
	valKeys := [2]integration.TestValidatorKey{keys[0], keys[1]}
	localCommitInfo, err := buildLocalCommitInfo(
		ctx,
		exchangeRates,
		valKeys,
		app.ChainID(),
	)
	s.Require().NoError(err)

	// update block header info and commit info
	ctx = buildCtxHeaderAndCommitInfo(ctx, localCommitInfo, app.ChainID())

	testCases := []struct {
		name                   string
		logger                 log.Logger
		oracleKeeper           oraclekeeper.Keeper
		stakingKeeper          *stakingkeeper.Keeper
		prepareProposalRequest *cometabci.RequestPrepareProposal
		expErr                 bool
		expErrMsg              string
	}{
		{
			name:          "nil prepare proposal request",
			logger:        log.NewNopLogger(),
			oracleKeeper:  app.OracleKeeper,
			stakingKeeper: app.StakingKeeper,
			expErr:        true,
			expErrMsg:     "prepare proposal received a nil request",
		},
		{
			name:          "prepare proposal request with nil Txs",
			logger:        log.NewNopLogger(),
			oracleKeeper:  app.OracleKeeper,
			stakingKeeper: app.StakingKeeper,
			prepareProposalRequest: &cometabci.RequestPrepareProposal{
				Height:          ctx.BlockHeight(),
				LocalLastCommit: localCommitInfo,
			},
			expErr:    true,
			expErrMsg: "prepare proposal received a request with nil Txs",
		},
		{
			name:          "prepare proposal request successful",
			logger:        log.NewNopLogger(),
			oracleKeeper:  app.OracleKeeper,
			stakingKeeper: app.StakingKeeper,
			prepareProposalRequest: &cometabci.RequestPrepareProposal{
				Height:          ctx.BlockHeight(),
				Txs:             [][]byte{},
				LocalLastCommit: localCommitInfo,
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			h := abci.NewProposalHandler(
				tc.logger,
				tc.oracleKeeper,
				tc.stakingKeeper,
			)

			resp, err := h.PrepareProposalHandler()(ctx, tc.prepareProposalRequest)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)

				var injectedVoteExtTx oracletypes.InjectedVoteExtensionTx
				err = injectedVoteExtTx.Unmarshal(resp.Txs[0])
				s.Require().NoError(err)

				sort.Slice(valKeys[:], func(i, j int) bool {
					return valKeys[i].ValAddress.String() < valKeys[j].ValAddress.String()
				})
				s.Require().Equal(exchangeRates, injectedVoteExtTx.ExchangeRateVotes[0].ExchangeRates)
				s.Require().Equal(valKeys[0].ValAddress.String(), injectedVoteExtTx.ExchangeRateVotes[0].Voter)
				s.Require().Equal(exchangeRates, injectedVoteExtTx.ExchangeRateVotes[1].ExchangeRates)
				s.Require().Equal(valKeys[1].ValAddress.String(), injectedVoteExtTx.ExchangeRateVotes[1].Voter)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestProcessProposalHandler() {
	app, ctx, keys := s.app, s.ctx, s.keys

	// enable vote extensions
	ctx = ctx.WithBlockHeight(3)
	consensusParams := ctx.ConsensusParams()
	consensusParams.Abci = &cmtproto.ABCIParams{
		VoteExtensionsEnableHeight: 2,
	}
	ctx = ctx.WithConsensusParams(consensusParams)

	// build local commit info
	exchangeRates := sdk.NewDecCoinsFromCoins(sdk.NewCoin("ATOM", math.NewInt(11)))
	valKeys := [2]integration.TestValidatorKey{keys[0], keys[1]}
	localCommitInfo, err := buildLocalCommitInfo(
		ctx,
		exchangeRates,
		valKeys,
		app.ChainID(),
	)
	s.Require().NoError(err)

	// build local commit info with different exchange rate
	exchangeRatesConflicting := sdk.NewDecCoinsFromCoins(sdk.NewCoin("ATOM", math.NewInt(111)))
	localCommitInfoConflicting, err := buildLocalCommitInfo(
		ctx,
		exchangeRatesConflicting,
		valKeys,
		app.ChainID(),
	)
	s.Require().NoError(err)

	// build injected vote extention tx
	sort.Slice(valKeys[:], func(i, j int) bool {
		return valKeys[i].ValAddress.String() < valKeys[j].ValAddress.String()
	})
	exchangeRateVotes := []oracletypes.AggregateExchangeRateVote{
		{
			ExchangeRates: exchangeRates,
			Voter:         valKeys[0].ValAddress.String(),
		},
		{
			ExchangeRates: exchangeRates,
			Voter:         valKeys[1].ValAddress.String(),
		},
	}
	localCommitInfoBz, err := localCommitInfo.Marshal()
	s.Require().NoError(err)
	injectedVoteExtTx := oracletypes.InjectedVoteExtensionTx{
		ExchangeRateVotes:  exchangeRateVotes,
		ExtendedCommitInfo: localCommitInfoBz,
	}
	bz, err := injectedVoteExtTx.Marshal()
	s.Require().NoError(err)
	var txs [][]byte
	txs = append(txs, bz)

	// create tx with conflicting local commit info
	localCommitInfoConflictingBz, err := localCommitInfoConflicting.Marshal()
	s.Require().NoError(err)
	injectedVoteExtTxConflicting := oracletypes.InjectedVoteExtensionTx{
		ExchangeRateVotes:  exchangeRateVotes,
		ExtendedCommitInfo: localCommitInfoConflictingBz,
	}
	bz, err = injectedVoteExtTxConflicting.Marshal()
	s.Require().NoError(err)
	var txsConflicting [][]byte
	txsConflicting = append(txsConflicting, bz)

	// update block header info and commit info
	ctx = buildCtxHeaderAndCommitInfo(ctx, localCommitInfo, app.ChainID())

	testCases := []struct {
		name                   string
		logger                 log.Logger
		oracleKeeper           oraclekeeper.Keeper
		stakingKeeper          *stakingkeeper.Keeper
		processProposalRequest *cometabci.RequestProcessProposal
		expErr                 bool
		expErrMsg              string
	}{
		{
			name:          "nil process proposal request",
			logger:        log.NewNopLogger(),
			oracleKeeper:  app.OracleKeeper,
			stakingKeeper: app.StakingKeeper,
			expErr:        true,
			expErrMsg:     "process proposal received a nil request",
		},
		{
			name:          "process proposal request with nil Txs",
			logger:        log.NewNopLogger(),
			oracleKeeper:  app.OracleKeeper,
			stakingKeeper: app.StakingKeeper,
			processProposalRequest: &cometabci.RequestProcessProposal{
				Height: ctx.BlockHeight(),
			},
			expErr:    true,
			expErrMsg: "process proposal received a request with nil Txs",
		},
		{
			name:          "process proposal request successful",
			logger:        log.NewNopLogger(),
			oracleKeeper:  app.OracleKeeper,
			stakingKeeper: app.StakingKeeper,
			processProposalRequest: &cometabci.RequestProcessProposal{
				Height: ctx.BlockHeight(),
				Txs:    txs,
			},
			expErr: false,
		},
		{
			name:          "process proposal request fails to verify exchange rate votes",
			logger:        log.NewNopLogger(),
			oracleKeeper:  app.OracleKeeper,
			stakingKeeper: app.StakingKeeper,
			processProposalRequest: &cometabci.RequestProcessProposal{
				Height: ctx.BlockHeight(),
				Txs:    txsConflicting,
			},
			expErr:    true,
			expErrMsg: "injected exhange rate votes and generated exchange votes are not equal",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			h := abci.NewProposalHandler(
				tc.logger,
				tc.oracleKeeper,
				tc.stakingKeeper,
			)

			resp, err := h.ProcessProposalHandler()(ctx, tc.processProposalRequest)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				s.Require().Equal(cometabci.ResponseProcessProposal_ACCEPT, resp.Status)
			}
		})
	}
}

func buildCtxHeaderAndCommitInfo(
	ctx sdk.Context,
	localCommitInfo cometabci.ExtendedCommitInfo,
	chainID string,
) sdk.Context {
	headerInfo := header.Info{
		Height:  ctx.BlockHeight(),
		Time:    ctx.BlockTime(),
		ChainID: chainID,
	}
	ctx = ctx.WithHeaderInfo(headerInfo)
	misbehavior := make([]cometabci.Misbehavior, 0)
	validatorHash := make([]byte, 0)
	proposerAddress := make([]byte, 0)
	lastCommit := cometabci.CommitInfo{
		Round: 0,
		Votes: []cometabci.VoteInfo{
			{
				Validator:   localCommitInfo.Votes[0].Validator,
				BlockIdFlag: cmtproto.BlockIDFlagCommit,
			},
			{
				Validator:   localCommitInfo.Votes[1].Validator,
				BlockIdFlag: cmtproto.BlockIDFlagCommit,
			},
		},
	}

	return ctx.WithCometInfo(baseapp.NewBlockInfo(misbehavior, validatorHash, proposerAddress, lastCommit))
}

// Builds local commit info with exchange rates and 2 validators
func buildLocalCommitInfo(
	ctx sdk.Context,
	exchangeRates sdk.DecCoins,
	valKeys [2]integration.TestValidatorKey,
	chainID string,
) (cometabci.ExtendedCommitInfo, error) {
	voteExt := oracletypes.OracleVoteExtension{
		ExchangeRates: exchangeRates,
		Height:        ctx.BlockHeight(),
	}
	bz, err := voteExt.Marshal()
	if err != nil {
		return cometabci.ExtendedCommitInfo{}, err
	}

	marshalDelimitedFn := func(msg proto.Message) ([]byte, error) {
		var buf bytes.Buffer
		if err := protoio.NewDelimitedWriter(&buf).WriteMsg(msg); err != nil {
			return nil, err
		}

		return buf.Bytes(), nil
	}
	cve := cmtproto.CanonicalVoteExtension{
		Extension: bz,
		Height:    ctx.BlockHeight() - 1, // the vote extension was signed in the previous height
		Round:     0,
		ChainId:   chainID,
	}
	extSignBytes, err := marshalDelimitedFn(&cve)
	if err != nil {
		return cometabci.ExtendedCommitInfo{}, err
	}

	votes := make([]cometabci.ExtendedVoteInfo, 2)
	for i := range votes {
		valConsAddr := sdk.ConsAddress(valKeys[i].PubKey.Address())
		extSig, err := valKeys[i].PrivKey.Sign(extSignBytes)
		if err != nil {
			return cometabci.ExtendedCommitInfo{}, err
		}

		votes[i] = cometabci.ExtendedVoteInfo{
			Validator: cometabci.Validator{
				Address: valConsAddr,
				Power:   100,
			},
			VoteExtension:      bz,
			BlockIdFlag:        cmtproto.BlockIDFlagCommit,
			ExtensionSignature: extSig,
		}
	}
	sort.Slice(votes, func(i, j int) bool {
		return string(votes[i].Validator.Address) < string(votes[j].Validator.Address)
	})

	return cometabci.ExtendedCommitInfo{
		Round: 0,
		Votes: votes,
	}, nil
}
