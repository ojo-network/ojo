package abci

import (
	"fmt"
	"sort"

	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	oraclekeeper "github.com/ojo-network/ojo/x/oracle/keeper"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

type ProposalHandler struct {
	logger        log.Logger
	oracleKeeper  oraclekeeper.Keeper
	stakingKeeper *stakingkeeper.Keeper
}

func NewProposalHandler(
	logger log.Logger,
	oracleKeeper oraclekeeper.Keeper,
	stakingKeeper *stakingkeeper.Keeper,
) *ProposalHandler {
	return &ProposalHandler{
		logger:        logger,
		oracleKeeper:  oracleKeeper,
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
			exchangeRateVotes, err := h.generateExchangeRateVotes(ctx, req.LocalLastCommit)
			if err != nil {
				return &cometabci.ResponsePrepareProposal{Txs: make([][]byte, 0)}, err
			}
			extendedCommitInfoBz, err := req.LocalLastCommit.Marshal()
			if err != nil {
				return &cometabci.ResponsePrepareProposal{Txs: make([][]byte, 0)}, err
			}

			medianGasEstimates, err := h.generateMedianGasEstimates(ctx, req.LocalLastCommit)
			if err != nil {
				return &cometabci.ResponsePrepareProposal{Txs: make([][]byte, 0)}, err
			}
			injectedVoteExtTx := oracletypes.InjectedVoteExtensionTx{
				ExchangeRateVotes:  exchangeRateVotes,
				ExtendedCommitInfo: extendedCommitInfoBz,
				GasEstimateMedians: medianGasEstimates,
			}

			bz, err := injectedVoteExtTx.Marshal()
			if err != nil {
				h.logger.Error("failed to encode injected vote extension tx", "err", err)
				return &cometabci.ResponsePrepareProposal{Txs: make([][]byte, 0)}, oracletypes.ErrEncodeInjVoteExt
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
		h.logger.Info("begin ProcessProposalHandler")
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

		h.logger.Info("after the first block")

		voteExtensionsEnabled := VoteExtensionsEnabled(ctx)
		if voteExtensionsEnabled {
			if len(req.Txs) < 1 {
				h.logger.Error("got process proposal request with no commit info")
				return &cometabci.ResponseProcessProposal{Status: cometabci.ResponseProcessProposal_REJECT},
					oracletypes.ErrNoCommitInfo
			}

			oracleRateFound := false
			var injectedVoteExtTx oracletypes.InjectedVoteExtensionTx
			for _, tx := range req.Txs {
				if err := injectedVoteExtTx.Unmarshal(tx); err == nil {
					oracleRateFound = true
					break
				}
			}
			if !(oracleRateFound) {
				h.logger.Error("failed to decode injected vote extension tx")
				return &cometabci.ResponseProcessProposal{
						Status: cometabci.ResponseProcessProposal_REJECT,
					},
					fmt.Errorf("failed to decode injected vote extension tx")
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
				h.logger.Error("failed to validate vote extensions", "err", err)
				return &cometabci.ResponseProcessProposal{Status: cometabci.ResponseProcessProposal_REJECT}, err
			}

			// Verify the proposer's oracle exchange rate votes by computing the same
			// calculation and comparing the results.
			exchangeRateVotes, err := h.generateExchangeRateVotes(ctx, extendedCommitInfo)
			if err != nil {
				h.logger.Error("failed to generate exchange rate votes", "err", err)
				fmt.Println("injected: ", injectedVoteExtTx.ExchangeRateVotes)
				fmt.Println("computed: ", exchangeRateVotes)
				return &cometabci.ResponseProcessProposal{Status: cometabci.ResponseProcessProposal_REJECT}, err
			}
			if err := h.verifyExchangeRateVotes(injectedVoteExtTx.ExchangeRateVotes, exchangeRateVotes); err != nil {
				h.logger.Error("failed to verify exchange rate votes", "err", err)
				fmt.Println("injected: ", injectedVoteExtTx.ExchangeRateVotes)
				fmt.Println("computed: ", exchangeRateVotes)
				return &cometabci.ResponseProcessProposal{Status: cometabci.ResponseProcessProposal_REJECT}, err
			}
			// Verify the proposer's gas estimation by computing the same median.
			gasEstimateMedians, err := h.generateMedianGasEstimates(ctx, extendedCommitInfo)
			if err != nil {
				h.logger.Error("failed to generate median gas estimates", "err", err)
				fmt.Println("injected: ", injectedVoteExtTx.GasEstimateMedians)
				fmt.Println("computed: ", gasEstimateMedians)
				return &cometabci.ResponseProcessProposal{Status: cometabci.ResponseProcessProposal_REJECT}, err
			}
			if err := h.verifyMedianGasEstimations(injectedVoteExtTx.GasEstimateMedians, gasEstimateMedians); err != nil {
				h.logger.Error("failed to verify median gas estimations", "err", err)
				fmt.Println("injected: ", injectedVoteExtTx.GasEstimateMedians)
				fmt.Println("computed: ", gasEstimateMedians)
				return &cometabci.ResponseProcessProposal{Status: cometabci.ResponseProcessProposal_REJECT}, err
			}
		}

		h.logger.Info(
			"processed proposal",
			"txs", len(req.Txs),
			"vote_extensions_enabled", voteExtensionsEnabled,
		)

		h.logger.Info("end ProcessProposalHandler")

		return &cometabci.ResponseProcessProposal{Status: cometabci.ResponseProcessProposal_ACCEPT}, nil
	}
}

func (h *ProposalHandler) generateExchangeRateVotes(
	ctx sdk.Context,
	ci cometabci.ExtendedCommitInfo,
) (votes []oracletypes.AggregateExchangeRateVote, err error) {
	for _, vote := range ci.Votes {
		if vote.BlockIdFlag != cmtproto.BlockIDFlagCommit {
			continue
		}

		var voteExt oracletypes.OracleVoteExtension
		if err := voteExt.Unmarshal(vote.VoteExtension); err != nil {
			h.logger.Error(
				"failed to decode vote extension",
				"err", err,
			)
			return nil, err
		}

		var valConsAddr sdk.ConsAddress
		if err := valConsAddr.Unmarshal(vote.Validator.Address); err != nil {
			h.logger.Error(
				"failed to unmarshal validator consensus address",
				"err", err,
			)
			return nil, err
		}
		val, err := h.stakingKeeper.GetValidatorByConsAddr(ctx, valConsAddr)
		if err != nil {
			h.logger.Error(
				"failed to get consensus validator from staking keeper",
				"err", err,
			)
			return nil, err
		}
		valAddr, err := sdk.ValAddressFromBech32(val.OperatorAddress)
		if err != nil {
			return nil, err
		}

		exchangeRateVote := oracletypes.NewAggregateExchangeRateVote(voteExt.ExchangeRates, valAddr)
		votes = append(votes, exchangeRateVote)
	}

	// sort votes so they are verified in the same order in ProcessProposalHandler
	sort.Slice(votes, func(i, j int) bool {
		return votes[i].Voter < votes[j].Voter
	})

	return votes, nil
}

func (h *ProposalHandler) verifyExchangeRateVotes(
	injectedVotes []oracletypes.AggregateExchangeRateVote,
	generatedVotes []oracletypes.AggregateExchangeRateVote,
) error {
	if len(injectedVotes) != len(generatedVotes) {
		return oracletypes.ErrNonEqualInjVotesLen
	}

	for i := range injectedVotes {
		injectedVote := injectedVotes[i]
		generatedVote := generatedVotes[i]

		if injectedVote.Voter != generatedVote.Voter || !injectedVote.ExchangeRates.Equal(generatedVote.ExchangeRates) {
			return oracletypes.ErrNonEqualInjVotesRates
		}
	}

	return nil
}

func (h *ProposalHandler) generateMedianGasEstimates(
	ctx sdk.Context,
	ci cometabci.ExtendedCommitInfo,
) ([]oracletypes.GasEstimate, error) {
	gasEstimates := []oracletypes.GasEstimate{}

	for _, vote := range ci.Votes {
		if vote.BlockIdFlag != cmtproto.BlockIDFlagCommit {
			continue
		}

		var valConsAddr sdk.ConsAddress
		if err := valConsAddr.Unmarshal(vote.Validator.Address); err != nil {
			h.logger.Error(
				"failed to unmarshal validator consensus address",
				"err", err,
			)
			return gasEstimates, err
		}
		val, err := h.stakingKeeper.GetValidatorByConsAddr(ctx, valConsAddr)
		if err != nil {
			h.logger.Error(
				"failed to get consensus validator from staking keeper",
				"err", err,
			)
			return gasEstimates, err
		}
		_, err = sdk.ValAddressFromBech32(val.OperatorAddress)
		if err != nil {
			return gasEstimates, err
		}
	}

	networks := []string{}
	// get contracts on registry list
	params := h.oracleKeeper.GasEstimateKeeper.GetParams(ctx)
	for _, contract := range params.ContractRegistry {
		networks = append(networks, contract.Network)
	}

	for _, network := range networks {
		networkEstimates := []oracletypes.GasEstimate{}

		for _, vote := range ci.Votes {
			var voteExt oracletypes.OracleVoteExtension
			if err := voteExt.Unmarshal(vote.VoteExtension); err != nil {
				h.logger.Error(
					"failed to decode vote extension",
					"err", err,
				)
				continue
			}

			for _, estimate := range voteExt.GasEstimates {
				if estimate.Network == network {
					networkEstimates = append(networkEstimates, estimate)
				}
			}
		}

		median, err := calculateMedian(networkEstimates)
		if err != nil {
			continue
		}
		gasEstimates = append(gasEstimates, median)
	}

	return gasEstimates, nil
}

func (h *ProposalHandler) verifyMedianGasEstimations(
	injectedEstimates []oracletypes.GasEstimate,
	generatedEstimates []oracletypes.GasEstimate,
) error {
	if len(injectedEstimates) != len(generatedEstimates) {
		return oracletypes.ErrNonEqualInjVotesLen
	}

	for i := range injectedEstimates {
		injectedEstimate := injectedEstimates[i]
		generatedEstimate := generatedEstimates[i]

		if injectedEstimate.Network != generatedEstimate.Network {
			return oracletypes.ErrNonEqualInjVotesRates
		}

		if injectedEstimate.GasEstimation != generatedEstimate.GasEstimation {
			return oracletypes.ErrNonEqualInjVotesRates
		}
	}

	return nil
}
