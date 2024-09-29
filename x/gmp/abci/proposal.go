package abci

import (
	"fmt"
	"sort"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	cometabci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	gmpkeeper "github.com/ojo-network/ojo/x/gmp/keeper"
	gmptypes "github.com/ojo-network/ojo/x/gmp/types"
)

type ProposalHandler struct {
	logger        log.Logger
	gmpKeeper     gmpkeeper.Keeper
	stakingKeeper stakingkeeper.Keeper
}

func NewProposalHandler(
	logger log.Logger,
	gmpKeeper gmpkeeper.Keeper,
	stakingKeeper stakingkeeper.Keeper,
) *ProposalHandler {
	return &ProposalHandler{
		logger:        logger,
		gmpKeeper:     gmpKeeper,
		stakingKeeper: stakingKeeper,
	}
}

// PrepareProposalHandler is called only on the selected validator as "block proposer" (selected by CometBFT, read
// more about this process here: https://docs.cometbft.com/v0.38/spec/consensus/proposer-selection). The block
// proposer is in charge of creating the next block by selecting the transactions from the mempool, and in this
// method it will create an extra transaction using the vote extension from the previous block which are only
// available on the next height at which vote extensions were enabled.
func (h *ProposalHandler) PrepareProposalHandler() sdk.PrepareProposalHandler {
	return func(ctx sdk.Context, req *cometabci.RequestPrepareProposal) (*cometabci.ResponsePrepareProposal, error) {
		if req == nil {
			err := fmt.Errorf("prepare proposal received a nil request")
			h.logger.Error(err.Error())
			return nil, err
		}

		err := baseapp.ValidateVoteExtensions(ctx, h.stakingKeeper, req.Height, ctx.ChainID(), req.LocalLastCommit)
		if err != nil {
			return &cometabci.ResponsePrepareProposal{Txs: make([][]byte, 0)}, err
		}

		if req.Txs == nil {
			err := fmt.Errorf("prepare proposal received a request with nil Txs")
			h.logger.Error(
				"height", req.Height,
				err.Error(),
			)
			return &cometabci.ResponsePrepareProposal{Txs: make([][]byte, 0)}, err
		}

		proposalTxs := req.Txs

		voteExtensionsEnabled := VoteExtensionsEnabled(ctx)
		if voteExtensionsEnabled {
			medianGasEstimation, err := h.generateMedianGasEstimate(ctx, req.LocalLastCommit)
			if err != nil {
				return &cometabci.ResponsePrepareProposal{Txs: make([][]byte, 0)}, err
			}
			extendedCommitInfoBz, err := req.LocalLastCommit.Marshal()
			if err != nil {
				return &cometabci.ResponsePrepareProposal{Txs: make([][]byte, 0)}, err
			}

			injectedVoteExtTx := gmptypes.InjectedVoteExtensionTx{
				MedianGasEstimation: &medianGasEstimation,
				ExtendedCommitInfo:  extendedCommitInfoBz,
			}

			bz, err := injectedVoteExtTx.Marshal()
			if err != nil {
				h.logger.Error("failed to encode injected vote extension tx", "err", err)
				return &cometabci.ResponsePrepareProposal{Txs: make([][]byte, 0)}, gmptypes.ErrEncodeInjVoteExt
			}

			// Inject a placeholder tx into the proposal s.t. validators can decode, verify,
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

// ProcessProposalHandler is called on all validators, and they can verify if the proposed block is valid. In case an
// invalid block is being proposed validators can reject it, causing a new round of PrepareProposal to happen. This
// step MUST be deterministic.
func (h *ProposalHandler) ProcessProposalHandler() sdk.ProcessProposalHandler {
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
			return &cometabci.ResponseProcessProposal{Status: cometabci.ResponseProcessProposal_REJECT}, err
		}

		voteExtensionsEnabled := VoteExtensionsEnabled(ctx)
		if voteExtensionsEnabled {
			if len(req.Txs) < 1 {
				h.logger.Error("got process proposal request with no commit info")
				return &cometabci.ResponseProcessProposal{Status: cometabci.ResponseProcessProposal_REJECT},
					gmptypes.ErrNoCommitInfo
			}

			var injectedVoteExtTx gmptypes.InjectedVoteExtensionTx
			if err := injectedVoteExtTx.Unmarshal(req.Txs[0]); err != nil {
				h.logger.Error("failed to decode injected vote extension tx", "err", err)
				return &cometabci.ResponseProcessProposal{Status: cometabci.ResponseProcessProposal_REJECT}, err
			}
			var extendedCommitInfo cometabci.ExtendedCommitInfo
			if err := extendedCommitInfo.Unmarshal(injectedVoteExtTx.ExtendedCommitInfo); err != nil {
				h.logger.Error("failed to decode injected extended commit info", "err", err)
				return &cometabci.ResponseProcessProposal{Status: cometabci.ResponseProcessProposal_REJECT}, err
			}

			err := baseapp.ValidateVoteExtensions(
				ctx,
				h.stakingKeeper,
				req.Height,
				ctx.ChainID(),
				extendedCommitInfo,
			)
			if err != nil {
				return &cometabci.ResponseProcessProposal{Status: cometabci.ResponseProcessProposal_REJECT}, err
			}

			// Verify the proposer's gas estimation by computing the same median.
			gasEstimateMedian, err := h.generateMedianGasEstimate(ctx, extendedCommitInfo)
			if err != nil {
				return &cometabci.ResponseProcessProposal{Status: cometabci.ResponseProcessProposal_REJECT}, err
			}
			if err := h.verifyMedianGasEstimation(*injectedVoteExtTx.MedianGasEstimation, gasEstimateMedian); err != nil {
				return &cometabci.ResponseProcessProposal{Status: cometabci.ResponseProcessProposal_REJECT}, err
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

func (h *ProposalHandler) generateMedianGasEstimate(
	ctx sdk.Context,
	ci cometabci.ExtendedCommitInfo,
) (median math.LegacyDec, err error) {
	gasEstimates := make([]math.LegacyDec, 0)
	for _, vote := range ci.Votes {
		if vote.BlockIdFlag != cmtproto.BlockIDFlagCommit {
			continue
		}

		var voteExt gmptypes.GmpVoteExtension
		if err := voteExt.Unmarshal(vote.VoteExtension); err != nil {
			h.logger.Error(
				"failed to decode vote extension",
				"err", err,
			)
			return math.LegacyZeroDec(), err
		}

		var valConsAddr sdk.ConsAddress
		if err := valConsAddr.Unmarshal(vote.Validator.Address); err != nil {
			h.logger.Error(
				"failed to unmarshal validator consensus address",
				"err", err,
			)
			return math.LegacyZeroDec(), err
		}
		val, err := h.stakingKeeper.GetValidatorByConsAddr(ctx, valConsAddr)
		if err != nil {
			h.logger.Error(
				"failed to get consensus validator from staking keeper",
				"err", err,
			)
			return math.LegacyZeroDec(), err
		}
		_, err = sdk.ValAddressFromBech32(val.OperatorAddress)
		if err != nil {
			return math.LegacyZeroDec(), err
		}

		// append median gas estimate to gas estimates
		gasEstimates = append(gasEstimates, *voteExt.GasEstimation)
	}

	// calculate median of gas estimates
	return calculateMedian(gasEstimates), nil
}

func (h *ProposalHandler) verifyMedianGasEstimation(
	injectedEstimation math.LegacyDec,
	generatedEstimation math.LegacyDec,
) error {
	// if they're not the same, error
	if injectedEstimation != generatedEstimation {
		return fmt.Errorf("injected median gas estimation does not match generated median gas estimation")
	}
	return nil
}

func calculateMedian(values []math.LegacyDec) math.LegacyDec {
	if len(values) == 0 {
		return math.LegacyZeroDec()
	}

	// Create a copy of the slice to avoid modifying the original
	sortedValues := make([]math.LegacyDec, len(values))
	copy(sortedValues, values)

	// Sort the copy in ascending order
	sort.Slice(sortedValues, func(i, j int) bool {
		return sortedValues[i].LT(sortedValues[j])
	})

	length := len(sortedValues)
	mid := length / 2

	if length%2 == 0 {
		// If even, return the average of the two middle values
		return sortedValues[mid-1].Add(sortedValues[mid]).QuoInt64(2)
	} else {
		// If odd, return the middle value
		return sortedValues[mid]
	}
}
