package abci

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/pricefeeder"
	"github.com/ojo-network/ojo/x/oracle/keeper"
	"github.com/ojo-network/ojo/x/oracle/types"
	"github.com/ojo-network/price-feeder/oracle"
)

// OracleVoteExtension defines the canonical vote extension structure.
type OracleVoteExtension struct {
	Height        int64
	ExchangeRates sdk.DecCoins
}

type VoteExtensionHandler struct {
	logger       log.Logger
	oracleKeeper keeper.Keeper
	priceFeeder  *pricefeeder.PriceFeeder
}

// NewVoteExtensionHandler returns a new VoteExtensionHandler.
func NewVoteExtensionHandler(
	logger log.Logger,
	oracleKeeper keeper.Keeper,
	priceFeeder *pricefeeder.PriceFeeder,
) *VoteExtensionHandler {
	return &VoteExtensionHandler{
		logger:       logger,
		oracleKeeper: oracleKeeper,
		priceFeeder:  priceFeeder,
	}
}

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

		// Get prices from Oracle Keeper's pricefeeder and generate vote msg
		if h.priceFeeder.Oracle == nil {
			err := fmt.Errorf("price feeder oracle not set")
			h.logger.Error(err.Error())
			return nil, err
		}
		prices := h.priceFeeder.Oracle.GetPrices()
		exchangeRatesStr := oracle.GenerateExchangeRatesString(prices)

		// Parse as DecCoins
		exchangeRates, err := types.ParseExchangeRateDecCoins(exchangeRatesStr)
		if err != nil {
			err := fmt.Errorf("extend vote handler received invalid exchange rate %w", types.ErrInvalidExchangeRate)
			h.logger.Error(
				"height", req.Height,
				err.Error(),
			)
			return nil, err
		}

		// Filter out rates which aren't included in the AcceptList
		acceptList := h.oracleKeeper.AcceptList(ctx)
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

		h.logger.Info(
			"verfied vote extension",
			"height", req.Height,
		)

		return &cometabci.ResponseVerifyVoteExtension{Status: cometabci.ResponseVerifyVoteExtension_ACCEPT}, nil
	}
}
