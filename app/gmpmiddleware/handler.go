package gmpmiddleware

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ojo-network/ojo/x/gmp/types"
)

const (
	blocksPerDay = 14400
)

type GmpKeeper interface {
	GetParams(ctx sdk.Context) (params types.Params)
	CreatePayment(
		goCtx context.Context,
		msg *types.MsgCreatePayment,
	) (*types.MsgCreatePaymentResponse, error)
}

type GmpHandler struct {
	gmp     GmpKeeper
	relayer string
}

func NewGmpHandler(k GmpKeeper, relayer string) *GmpHandler {
	return &GmpHandler{
		gmp:     k,
		relayer: relayer,
	}
}

// HandleGeneralMessage takes the receiving message from axelar,
// and sends it along to the GMP module.
func (h GmpHandler) HandleGeneralMessage(
	ctx sdk.Context,
	srcChain,
	srcAddress string,
	receiver string,
	payload []byte,
	sender string,
	channel string,
	coin sdk.Coin,
) error {
	ctx.Logger().Info("HandleGeneralMessage called",
		"srcChain", srcChain,
		"srcAddress", srcAddress, // this is the Ojo contract address
		"receiver", receiver,
		"payload", payload,
		"coin", coin,
	)

	err := verifyParams(h.gmp.GetParams(ctx), sender, channel)
	if err != nil {
		return err
	}
	msg, err := types.NewGmpDecoder(payload)
	if err != nil {
		return err
	}
	ctx.Logger().Info("HandleGeneralMessage GMP Decoder", "msg", msg)
	denoms := msg.GetDenoms()

	defaultDeviation, err := math.LegacyNewDecFromStr("1")
	if err != nil {
		return err
	}

	tx := &types.MsgCreatePayment{
		Relayer: h.relayer,
		Payment: &types.Payment{
			Relayer:          h.relayer,
			DestinationChain: srcChain,
			Denom:            denoms[0],
			Token:            coin,
			Deviation:        defaultDeviation,
			Heartbeat:        blocksPerDay,
		},
	}
	err = tx.ValidateBasic()
	if err != nil {
		return err
	}
	ctx.Logger().Info("HandleGeneralMessage GMP Decoder", "tx", tx)
	_, err = h.gmp.CreatePayment(ctx, tx)
	return err
}
