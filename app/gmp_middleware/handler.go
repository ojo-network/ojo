package gmp_middleware

import (
	"context"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gmptypes "github.com/ojo-network/ojo/x/gmp/types"
)

type GmpKeeper interface {
	RelayPrice(goCtx context.Context, msg *gmptypes.MsgRelayPrice) (*gmptypes.MsgRelayPriceResponse, error)
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
func (h GmpHandler) HandleGeneralMessage(ctx sdk.Context, srcChain, srcAddress string, destAddress string, payload []byte) error {
	ctx.Logger().Info("HandleGeneralMessage called",
		"srcChain", srcChain,
		"srcAddress", srcAddress,
		"destAddress", destAddress,
		"payload", payload,
		"module", "x/gmp-middleware",
	)

	denomString := string(payload)
	denoms := strings.Split(denomString, ",")

	_, err := h.gmp.RelayPrice(ctx, &gmptypes.MsgRelayPrice{
		Relayer:            srcAddress,
		DestinationChain:   srcChain,
		DestinationAddress: destAddress,
		Denoms:             denoms,
	},
	)

	return err
}

// HandleGeneralMessageWithToken currently performs a no-op.
func (h GmpHandler) HandleGeneralMessageWithToken(ctx sdk.Context, srcChain, srcAddress string, destAddress string, payload []byte, coin sdk.Coin) error {
	ctx.Logger().Info("HandleGeneralMessageWithToken called",
		"srcChain", srcChain,
		"srcAddress", srcAddress,
		"destAddress", destAddress,
		"payload", payload,
		"coin", coin,
	)

	return nil
}
