package abci

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

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

// PreBlocker is run before finalize block to update the aggregrate exchange rate votes on the oracle module
// that were verified by the vote etension handler so that the exchange rate votes are available during the
// entire block execution (from BeginBlock). It takes the module manger from app.go to execute the PreBlock
// methods of the other modules set in SetOrderPreBlockers.
func (h *PreBlockHandler) PreBlocker(mm *module.Manager) sdk.PreBlocker {
	return func(ctx sdk.Context, req *cometabci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
		if req == nil {
			err := fmt.Errorf("preblocker received a nil request")
			h.logger.Error(err.Error())
			return nil, err
		}

		// execute preblockers of modules in OrderPreBlockers first.
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		paramsChanged := false
		for _, moduleName := range mm.OrderPreBlockers {
			if module, ok := mm.Modules[moduleName].(appmodule.HasPreBlocker); ok {
				rsp, err := module.PreBlock(ctx)
				if err != nil {
					return nil, err
				}
				if rsp.IsConsensusParamsChanged() {
					paramsChanged = true
				}
			}
		}

		res := &sdk.ResponsePreBlock{
			ConsensusParamsChanged: paramsChanged,
		}

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
					h.logger.Error("failed to get voter address", "err", err)
					continue
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
