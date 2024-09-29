package abci

import (
	"fmt"

	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	relayerClient "github.com/ojo-network/ojo-evm/relayer/relayer/client"
	"github.com/ojo-network/ojo/x/gmp/keeper"
	"github.com/ojo-network/ojo/x/gmp/types"
)

type VoteExtensionHandler struct {
	logger    log.Logger
	gmpKeeper keeper.Keeper
}

// NewVoteExtensionHandler returns a new VoteExtensionHandler.
func NewVoteExtensionHandler(
	logger log.Logger,
	gmpKeeper keeper.Keeper,
) *VoteExtensionHandler {
	return &VoteExtensionHandler{
		logger:    logger,
		gmpKeeper: gmpKeeper,
	}
}

// ExtendVoteHandler creates an GmpVoteExtension using the prices fetched from the price feeder
// service. It will filter out exchange rates that are not part of the oracle module's accept list.
func (h *VoteExtensionHandler) ExtendVoteHandler() sdk.ExtendVoteHandler {
	return func(ctx sdk.Context, req *cometabci.RequestExtendVote) (resp *cometabci.ResponseExtendVote, err error) {
		defer func() {
			// catch panics if possible
			if r := recover(); r != nil {
				h.logger.Error(
					"recovered from panic in ExtendVoteHandler",
					"err", r,
				)

				resp, err = &cometabci.ResponseExtendVote{VoteExtension: []byte{}},
					fmt.Errorf("recovered application panic in ExtendVote: %v", r)
			}
		}()

		if req == nil {
			err := fmt.Errorf("extend vote handler received a nil request")
			h.logger.Error(err.Error())
			return nil, err
		}

		// get the current gas estimation for each contract
		gasEstimate, err := relayerClient.EstimateGasFee("Ethereum", "0x0000000000000000000000000000000000000000", "1000000", "1.5")
		if err != nil {
			h.logger.Error("error estimating gas fee", "error", err)
			return &cometabci.ResponseExtendVote{VoteExtension: []byte{}}, err
		}
		gasEstimateDec := gasEstimate.ToLegacyDec()

		voteExt := types.GmpVoteExtension{
			Height:        req.Height,
			GasEstimation: &gasEstimateDec,
		}

		bz, err := voteExt.Marshal()
		if err != nil {
			err := fmt.Errorf("failed to marshal vote extension: %w", err)
			h.logger.Error(
				"height", req.Height,
				err.Error(),
			)
			return &cometabci.ResponseExtendVote{VoteExtension: []byte{}}, err
		}

		h.logger.Info(
			"created vote extension",
			"height", req.Height,
		)

		return &cometabci.ResponseExtendVote{VoteExtension: bz}, nil
	}
}

// VerifyVoteExtensionHandler validates the OracleVoteExtension created by the ExtendVoteHandler. It
// verifies that the vote extension can unmarshal correctly and is for the correct height.
func (h *VoteExtensionHandler) VerifyVoteExtensionHandler() sdk.VerifyVoteExtensionHandler {
	return func(ctx sdk.Context, req *cometabci.RequestVerifyVoteExtension) (
		*cometabci.ResponseVerifyVoteExtension,
		error,
	) {
		if req == nil {
			err := fmt.Errorf("verify vote extension handler received a nil request")
			h.logger.Error(err.Error())
			return nil, err
		}

		if len(req.VoteExtension) == 0 {
			h.logger.Info(
				"verify vote extension handler received empty vote extension",
				"height", req.Height,
			)

			return &cometabci.ResponseVerifyVoteExtension{Status: cometabci.ResponseVerifyVoteExtension_ACCEPT}, nil
		}

		var voteExt types.GmpVoteExtension
		err := voteExt.Unmarshal(req.VoteExtension)
		if err != nil {
			err := fmt.Errorf("verify vote extension handler failed to unmarshal vote extension: %w", err)
			h.logger.Error(
				"height", req.Height,
				err.Error(),
			)
			return &cometabci.ResponseVerifyVoteExtension{Status: cometabci.ResponseVerifyVoteExtension_REJECT}, err
		}

		if voteExt.Height != req.Height {
			err := fmt.Errorf(
				"verify vote extension handler received vote extension height that doesn't"+
					"match request height; expected: %d, got: %d",
				req.Height,
				voteExt.Height,
			)
			h.logger.Error(err.Error())
			return &cometabci.ResponseVerifyVoteExtension{Status: cometabci.ResponseVerifyVoteExtension_REJECT}, err
		}

		h.logger.Info(
			"verfied vote extension",
			"height", req.Height,
		)

		return &cometabci.ResponseVerifyVoteExtension{Status: cometabci.ResponseVerifyVoteExtension_ACCEPT}, nil
	}
}
