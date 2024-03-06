package abci

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/x/oracle/keeper"
)

type PreBlockHandler struct {
	logger log.Logger
	keeper keeper.Keeper
}

func NewPreBlockHandler(logger log.Logger, keeper keeper.Keeper) *PreBlockHandler {
	return &PreBlockHandler{
		logger: logger,
		keeper: keeper,
	}
}

func (h *PreBlockHandler) PreBlocker() sdk.PreBlocker {
	return func(ctx sdk.Context, req *cometabci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
		if req == nil {
			err := fmt.Errorf("preblocker received a nil request")
			h.logger.Error(err.Error())
			return nil, err
		}

		res := &sdk.ResponsePreBlock{}
		if len(req.Txs) == 0 {
			return res, nil
		}

		voteExtensionsEnabled := VoteExtensionsEnabled(ctx)
		if voteExtensionsEnabled {
			var injectedVoteExtTx AggregateExchangeRateVotes
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

		h.logger.Info(
			"oracle preblocker executed",
			"vote_extensions_enabled", voteExtensionsEnabled,
		)

		return res, nil
	}
}
