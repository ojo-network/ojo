package gmpmiddleware

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ojo-network/ojo/x/gmp/types"
)

type GmpKeeper interface {
	RelayPrice(
		goCtx context.Context,
		msg *types.MsgRelayPrice,
	) (*types.MsgRelayPriceResponse, error)
	GetParams(ctx sdk.Context) (params types.Params)
}

type GmpHandler struct {
	gmp GmpKeeper
}

func NewGmpHandler(k GmpKeeper) *GmpHandler {
	return &GmpHandler{
		gmp: k,
	}
}

// HandleGeneralMessage takes the receiving message from axelar,
// and sends it along to the GMP module.
func (h GmpHandler) HandleGeneralMessage(
	ctx sdk.Context,
	srcChain,
	srcAddress string,
	destAddress string,
	payload []byte,
	sender string,
	channel string,
) error {
	ctx.Logger().Info("HandleGeneralMessage called",
		"srcChain", srcChain,
		"srcAddress", srcAddress,
		"destAddress", destAddress,
		"payload", payload,
		"module", "x/gmp-middleware",
	)

	params := h.gmp.GetParams(ctx)
	if !strings.EqualFold(params.GmpAddress, sender) {
		return fmt.Errorf("invalid sender address: %s", sender)
	}
	if !strings.EqualFold(params.GmpChannel, channel) {
		return fmt.Errorf("invalid channel: %s", channel)
	}

	denomString := string(payload)
	denoms := strings.Split(denomString, ",")

	_, err := h.gmp.RelayPrice(ctx, &types.MsgRelayPrice{
		Relayer:            srcAddress,
		DestinationChain:   srcChain,
		DestinationAddress: destAddress,
		Denoms:             denoms,
	},
	)

	return err
}

// HandleGeneralMessage takes the receiving message from axelar,
// and sends it along to the GMP module.
func (h GmpHandler) HandleGeneralMessageWithToken(
	ctx sdk.Context,
	srcChain,
	srcAddress string,
	destAddress string,
	payload []byte,
	sender string,
	channel string,
	coin sdk.Coin,
) error {
	ctx.Logger().Info("HandleGeneralMessageWithToken called",
		"srcChain", srcChain,
		"srcAddress", srcAddress,
		"destAddress", destAddress,
		"payload", payload,
		"coin", coin,
	)

	params := h.gmp.GetParams(ctx)
	if !strings.EqualFold(params.GmpAddress, sender) {
		return fmt.Errorf("invalid sender address: %s", sender)
	}
	if !strings.EqualFold(params.GmpChannel, channel) {
		return fmt.Errorf("invalid channel: %s", channel)
	}

	denomString := string(payload)
	denoms := strings.Split(denomString, ",")

	_, err := h.gmp.RelayPrice(ctx, &types.MsgRelayPrice{
		Relayer:            srcAddress,
		DestinationChain:   srcChain,
		DestinationAddress: destAddress,
		Denoms:             denoms,
		Token:              coin,
	},
	)

	return err
}
