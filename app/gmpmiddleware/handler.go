package gmpmiddleware

import (
	"context"

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
	tx := &types.MsgRelayPrice{
		Relayer:               h.relayer,
		DestinationChain:      srcChain,
		ClientContractAddress: msg.ContractAddress.Hex(),
		OjoContractAddress:    srcAddress,
		Denoms:                msg.GetDenoms(),
		CommandSelector:       msg.CommandSelector[:],
		CommandParams:         msg.CommandParams,
		Timestamp:             msg.Timestamp.Int64(),
		Token:                 coin,
	}
	err = tx.ValidateBasic()
	if err != nil {
		return err
	}
	ctx.Logger().Info("HandleGeneralMessage GMP Decoder", "tx", tx)
	_, err = h.gmp.RelayPrice(ctx, tx)
	return err
}
