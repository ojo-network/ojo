package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ojo-network/ojo/x/gmp/types"
)

type msgServer struct {
	keeper Keeper
}

// NewMsgServerImpl returns an implementation of the gmp MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{keeper: keeper}
}

// SetParams implements MsgServer.SetParams method.
// It defines a method to update the x/gmp module parameters.
func (ms msgServer) SetParams(goCtx context.Context, msg *types.MsgSetParams) (*types.MsgSetParamsResponse, error) {
	if ms.keeper.authority != msg.Authority {
		err := errors.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority; expected %s, got %s",
			ms.keeper.authority,
			msg.Authority,
		)
		return nil, err
	}

	if err := msg.Params.Validate(); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	ms.keeper.SetParams(ctx, *msg.Params)

	return &types.MsgSetParamsResponse{}, nil
}

// Relay implements MsgServer.Relay method.
// It defines a method to relay over GMP to recipient chains.
func (ms msgServer) RelayPrice(
	goCtx context.Context,
	msg *types.MsgRelayPrice,
) (*types.MsgRelayPriceResponse, error) {
	return ms.keeper.RelayPrice(goCtx, msg)
}

func (ms msgServer) CreatePayment(
	goCtx context.Context,
	msg *types.MsgCreatePayment,
) (*types.MsgCreatePaymentResponse, error) {
	return ms.keeper.CreatePayment(goCtx, msg)
}
