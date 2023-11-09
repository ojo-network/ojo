package keeper

import (
	"context"
	"math/big"
	"time"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
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
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := ms.keeper.GetParams(ctx)

	// encode oracle data
	rates := []types.PriceFeedData{}
	for _, denom := range msg.Denoms {
		rate, err := ms.keeper.oracleKeeper.GetExchangeRate(ctx, denom)
		if err != nil {
			return &types.MsgRelayPriceResponse{}, err
		}

		priceFeed, err := types.NewPriceFeedData(
			denom,
			rate,
			// TODO: replace with actual resolve time & id
			// Ref: https://github.com/ojo-network/ojo/issues/309
			big.NewInt(1),
			big.NewInt(1),
		)
		if err != nil {
			ms.keeper.Logger(ctx).With(err).Error("unable to relay price to gmp")
			continue
		}

		rates = append(rates, priceFeed)
	}

	// TODO: fill with actual disableResolve option
	// Ref: https://github.com/ojo-network/ojo/issues/309
	payload, err := types.EncodeABI("postPrices", rates, false)
	if err != nil {
		return nil, err
	}

	// package GMP
	message := types.GmpMessage{
		DestinationChain:   msg.DestinationChain,
		DestinationAddress: msg.DestinationAddress,
		Payload:            payload,
		Type:               types.TypeGeneralMessage,
	}
	bz, err := message.Marshal()
	if err != nil {
		return nil, err
	}

	// submit IBC transfer
	transferMsg := ibctransfertypes.NewMsgTransfer(
		ibctransfertypes.PortID,
		params.GmpChannel,
		msg.Token,
		msg.Relayer,
		params.GmpAddress,
		clienttypes.ZeroHeight(),
		uint64(ctx.BlockTime().Add(time.Duration(params.GmpTimeout)*time.Hour).UnixNano()),
		string(bz),
	)

	_, err = ms.keeper.ibcKeeper.Transfer(ctx, transferMsg)
	if err != nil {
		return &types.MsgRelayPriceResponse{}, err
	}

	return &types.MsgRelayPriceResponse{}, nil
}