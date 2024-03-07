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

type AggregateExchangeRateVotes struct {
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
		if req == nil {
			err := fmt.Errorf("prepare proposal received a nil request")
			h.logger.Error(err.Error())
			return nil, err
		}

		err := baseapp.ValidateVoteExtensions(ctx, h.valStore, req.Height, ctx.ChainID(), req.LocalLastCommit)
		if err != nil {
			return nil, err
		}

		if req.Txs == nil {
			err := fmt.Errorf("prepare proposal received a request with nil Txs")
			h.logger.Error(
				"height", req.Height,
				err.Error(),
			)
			return nil, err
		}

		proposalTxs := req.Txs

		voteExtensionsEnabled := VoteExtensionsEnabled(ctx)
		if voteExtensionsEnabled {
			exchangeRateVotes, err := h.generateExchangeRateVotes(req.LocalLastCommit)
			if err != nil {
				return nil, errors.New("failed to generate exchange rate votes")
			}

			injectedVoteExtTx := AggregateExchangeRateVotes{
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

		h.logger.Info(
			"prepared proposal",
			"txs", len(proposalTxs),
			"vote_extensions_enabled", voteExtensionsEnabled,
		)

		return &cometabci.ResponsePrepareProposal{
			Txs: proposalTxs,
		}, nil
	}
}

func (h *ProposalHandler) ProcessProposal() sdk.ProcessProposalHandler {
	return func(ctx sdk.Context, req *cometabci.RequestProcessProposal) (*cometabci.ResponseProcessProposal, error) {
		if req == nil {
			err := fmt.Errorf("process proposal received a nil request")
			h.logger.Error(err.Error())
			return nil, err
		}

		if req.Txs == nil {
			err := fmt.Errorf("process proposal received a request with nil Txs")
			h.logger.Error(
				"height", req.Height,
				err.Error(),
			)
			return nil, err
		}

		voteExtensionsEnabled := VoteExtensionsEnabled(ctx)
		if voteExtensionsEnabled {
			var injectedVoteExtTx AggregateExchangeRateVotes
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
			exchangeRateVotes, err := h.generateExchangeRateVotes(injectedVoteExtTx.ExtendedCommitInfo)
			if err != nil {
				return &cometabci.ResponseProcessProposal{Status: cometabci.ResponseProcessProposal_REJECT},
					errors.New("failed to generate exchange rate votes")
			}
			if len(injectedVoteExtTx.ExchangeRateVotes) != len(exchangeRateVotes) {
				return &cometabci.ResponseProcessProposal{Status: cometabci.ResponseProcessProposal_REJECT},
					errors.New("number of votes in vote extension and extended commit info are not equal")
			}
		}

		h.logger.Info(
			"processed proposal",
			"txs", len(req.Txs),
			"vote_extensions_enabled", voteExtensionsEnabled,
		)

		return &cometabci.ResponseProcessProposal{Status: cometabci.ResponseProcessProposal_ACCEPT}, nil
	}
}

func (h *ProposalHandler) generateExchangeRateVotes(
	ci cometabci.ExtendedCommitInfo,
) (votes []oracletypes.AggregateExchangeRateVote, err error) {
	for _, vote := range ci.Votes {
		if vote.BlockIdFlag != cmtproto.BlockIDFlagCommit {
			continue
		}

		var voteExt OracleVoteExtension
		if err := json.Unmarshal(vote.VoteExtension, &voteExt); err != nil {
			h.logger.Error(
				"failed to decode vote extension",
				"err", err,
				"validator", fmt.Sprintf("%x", vote.Validator.Address),
			)
			return nil, err
		}
		exchangeRateVote := oracletypes.NewAggregateExchangeRateVote(voteExt.ExchangeRates, voteExt.ValidatorAddress)
		votes = append(votes, exchangeRateVote)
	}

	return votes, nil
}
