package abci

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/x/oracle/keeper"
	"github.com/ojo-network/ojo/x/oracle/types"
	"github.com/ojo-network/price-feeder/oracle"
)

// OracleVoteExtension defines the canonical vote extension structure.
type OracleVoteExtension struct {
	Height        int64
	ExchangeRates sdk.DecCoins
}

type VoteExtHandler struct {
	logger       log.Logger
	OracleKeeper keeper.Keeper
}

// NewVoteExtensionHandler returns a new VoteExtensionHandler.
func NewVoteExtensionHandler(
	logger log.Logger,
	oracleKeeper keeper.Keeper,
) *VoteExtHandler {
	return &VoteExtHandler{
		logger:       logger,
		OracleKeeper: oracleKeeper,
	}
}

func (h *VoteExtHandler) ExtendVoteHandler() sdk.ExtendVoteHandler {
	return func(ctx sdk.Context, req *cometabci.RequestExtendVote) (*cometabci.ResponseExtendVote, error) {
		if req == nil {
			err := fmt.Errorf("extend vote handler received a nil request")
			h.logger.Error(
				"height", req.Height,
				err.Error(),
			)
			return nil, err
		}

		// Get prices from Oracle Keeper's pricefeeder and generate vote msg
		prices := h.OracleKeeper.PriceFeederOracle.GetPrices()
		h.logger.Info("Oracle price feeder prices", "prices", prices)
		exchangeRatesStr := oracle.GenerateExchangeRatesString(prices)

		// Parse as DecCoins
		exchangeRates, err := types.ParseExchangeRateDecCoins(exchangeRatesStr)
		if err != nil {
			err := fmt.Errorf("extend vote handler received invalid exchange rate", types.ErrInvalidExchangeRate)
			h.logger.Error(
				"height", req.Height,
				err.Error(),
			)
			return nil, err
		}
		// Filter out rates which aren't included in the AcceptList
		acceptList := h.OracleKeeper.AcceptList(ctx)
		filteredDecCoins := sdk.DecCoins{}
		for _, decCoin := range exchangeRates {
			if acceptList.Contains(decCoin.Denom) {
				filteredDecCoins = append(filteredDecCoins, decCoin)
			}
		}

		voteExt := OracleVoteExtension{
			Height:        req.Height,
			ExchangeRates: filteredDecCoins,
		}

		bz, err := json.Marshal(voteExt)
		if err != nil {
			err := fmt.Errorf("failed to marshal vote extension: %w", err)
			h.logger.Error(
				"height", req.Height,
				err.Error(),
			)
			return nil, err
		}
		h.logger.Info(
			"created vote extension",
			"height", req.Height,
		)

		return &cometabci.ResponseExtendVote{VoteExtension: bz}, nil
	}
}

func (h *VoteExtHandler) VerifyVoteExtensionHandler() sdk.VerifyVoteExtensionHandler {
	return func(ctx sdk.Context, req *cometabci.RequestVerifyVoteExtension) (*cometabci.ResponseVerifyVoteExtension, error) {
		if req == nil {
			err := fmt.Errorf("verify vote extension handler received a nil request")
			h.logger.Error(
				"height", req.Height,
				err.Error(),
			)
			return nil, err
		}

		var voteExt OracleVoteExtension
		err := json.Unmarshal(req.VoteExtension, &voteExt)
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

		// Verify vote extension signer and exchange rate vote signer match
		validatorAddress := sdk.ConsAddress{}
		if err := validatorAddress.Unmarshal(req.ValidatorAddress); err != nil {
			err := fmt.Errorf("verify vote extension handler failed to unmarshal validator address: %w", err)
			h.logger.Error(
				"height", req.Height,
				err.Error(),
			)
			return nil, err
		}

		h.logger.Info(
			"verfied vote extension",
			"height", req.Height,
		)

		return &cometabci.ResponseVerifyVoteExtension{Status: cometabci.ResponseVerifyVoteExtension_ACCEPT}, nil
	}
}
