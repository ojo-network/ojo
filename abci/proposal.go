package abci

import (
	"encoding/json"
	"errors"
	"fmt"

	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/x/oracle/keeper"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

type OracleExchangeRateVotes struct {
	ExchangeRateVotes  []oracletypes.AggregateExchangeRateVote
	ExtendedCommitInfo cometabci.ExtendedCommitInfo
}

type ProposalHandler struct {
	logger   log.Logger
	keeper   keeper.Keeper
	valStore baseapp.ValidatorStore
}

func NewProposalHandler(logger log.Logger, keeper keeper.Keeper, valStore baseapp.ValidatorStore) *ProposalHandler {
	return &ProposalHandler{
		logger:   logger,
		keeper:   keeper,
		valStore: valStore,
	}
}

func (h *ProposalHandler) PrepareProposal() sdk.PrepareProposalHandler {
	return func(ctx sdk.Context, req *cometabci.RequestPrepareProposal) (*cometabci.ResponsePrepareProposal, error) {
		err := baseapp.ValidateVoteExtensions(ctx, h.valStore, req.Height, ctx.ChainID(), req.LocalLastCommit)
		if err != nil {
			return nil, err
		}

		proposalTxs := req.Txs

		if req.Height > ctx.ConsensusParams().Abci.VoteExtensionsEnableHeight {
			exchangeRateVotes, err := h.generateExchangeRateVotes(ctx, req.LocalLastCommit)
			if err != nil {
				return nil, errors.New("failed to generate exchange rate votes")
			}

			injectedVoteExtTx := OracleExchangeRateVotes{
				ExchangeRateVotes:  exchangeRateVotes,
				ExtendedCommitInfo: req.LocalLastCommit,
			}

			// TODO: Switch from stdlib JSON encoding to a more performant mechanism.
			bz, err := json.Marshal(injectedVoteExtTx)
			if err != nil {
				h.logger.Error("failed to encode injected vote extension tx", "err", err)
				return nil, errors.New("failed to encode injected vote extension tx")
			}

			// Inject a "fake" tx into the proposal s.t. validators can decode, verify,
			// and store the oracle exchange rate votes.
			proposalTxs = append([][]byte{bz}, proposalTxs...)
		}

		return &cometabci.ResponsePrepareProposal{
			Txs: proposalTxs,
		}, nil
	}
}

func (h *ProposalHandler) ProcessProposal() sdk.ProcessProposalHandler {
	return func(ctx sdk.Context, req *cometabci.RequestProcessProposal) (*cometabci.ResponseProcessProposal, error) {
		if req.Height > ctx.ConsensusParams().Abci.VoteExtensionsEnableHeight {
			var injectedVoteExtTx OracleExchangeRateVotes
			if err := json.Unmarshal(req.Txs[0], &injectedVoteExtTx); err != nil {
				h.logger.Error("failed to decode injected vote extension tx", "err", err)
				return &cometabci.ResponseProcessProposal{Status: cometabci.ResponseProcessProposal_REJECT}, nil
			}

			err := baseapp.ValidateVoteExtensions(
				ctx,
				h.valStore,
				req.Height,
				ctx.ChainID(),
				injectedVoteExtTx.ExtendedCommitInfo,
			)
			if err != nil {
				return nil, err
			}

			// Verify the proposer's oracle exchange rate votes by computing the same
			// calculation and comparing the results.
		}

		return &cometabci.ResponseProcessProposal{Status: cometabci.ResponseProcessProposal_ACCEPT}, nil
	}
}

func (h *ProposalHandler) PreBlocker(ctx sdk.Context, req *cometabci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
	res := &sdk.ResponsePreBlock{}
	if len(req.Txs) == 0 {
		return res, nil
	}

	if req.Height > ctx.ConsensusParams().Abci.VoteExtensionsEnableHeight {
		var injectedVoteExtTx OracleExchangeRateVotes
		if err := json.Unmarshal(req.Txs[0], &injectedVoteExtTx); err != nil {
			h.logger.Error("failed to decode injected vote extension tx", "err", err)
			return nil, err
		}

		// set oracle exchange rate votes using the passed in context, which will make
		// these votes available in the current block.
		for _, exchangeRateVote := range injectedVoteExtTx.ExchangeRateVotes {
			valAddr, err := sdk.ValAddressFromBech32(exchangeRateVote.Voter)
			if err != nil {
				return nil, err
			}
			h.keeper.SetAggregateExchangeRateVote(ctx, valAddr, exchangeRateVote)
		}
	}
	return res, nil
}

func (h *ProposalHandler) generateExchangeRateVotes(
	ctx sdk.Context,
	ci cometabci.ExtendedCommitInfo,
) (votes []oracletypes.AggregateExchangeRateVote, err error) {
	for _, v := range ci.Votes {
		if v.BlockIdFlag != cmtproto.BlockIDFlagCommit {
			continue
		}

		var voteExt OracleVoteExtension
		if err := json.Unmarshal(v.VoteExtension, &voteExt); err != nil {
			h.logger.Error(
				"failed to decode vote extension",
				"err", err,
				"validator", fmt.Sprintf("%x", v.Validator.Address),
			)
			return nil, err
		}
		votes = append(votes, voteExt.ExchangeRateVote)
	}

	return votes, nil
}
