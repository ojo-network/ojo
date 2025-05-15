package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ojo-network/ojo/x/oracle/types"
)

func (ms msgServer) SetPriceFeeder(
	goCtx context.Context,
	msg *types.MsgSetPriceFeeder,
) (*types.MsgSetPriceFeederResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	feederAccount := sdk.MustAccAddressFromBech32(msg.Feeder)
	_, found := ms.Keeper.GetPriceFeeder(ctx, feederAccount)
	if !found {
		return nil, types.ErrNotAPriceFeeder
	}
	ms.Keeper.SetPriceFeeder(ctx, types.PriceFeeder{
		Feeder:   msg.Feeder,
		IsActive: msg.IsActive,
	})
	return &types.MsgSetPriceFeederResponse{}, nil
}

func (ms msgServer) DeletePriceFeeder(
	goCtx context.Context,
	msg *types.MsgDeletePriceFeeder,
) (*types.MsgDeletePriceFeederResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	feederAccount := sdk.MustAccAddressFromBech32(msg.Feeder)
	_, found := ms.Keeper.GetPriceFeeder(ctx, feederAccount)
	if !found {
		return nil, types.ErrNotAPriceFeeder
	}
	ms.RemovePriceFeeder(ctx, feederAccount)
	return &types.MsgDeletePriceFeederResponse{}, nil
}

func (ms msgServer) AddPriceFeeders(
	goCtx context.Context,
	msg *types.MsgAddPriceFeeders,
) (*types.MsgAddPriceFeedersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if ms.authority != msg.Authority {
		return nil, errors.Wrapf(
			types.ErrNoGovAuthority,
			"invalid authority; expected %s, got %s",
			ms.authority, msg.Authority,
		)
	}

	for _, feeder := range msg.Feeders {
		ms.Keeper.SetPriceFeeder(ctx, types.PriceFeeder{
			Feeder:   feeder,
			IsActive: true,
		})
	}
	return &types.MsgAddPriceFeedersResponse{}, nil
}

func (ms msgServer) RemovePriceFeeders(
	goCtx context.Context,
	msg *types.MsgRemovePriceFeeders,
) (*types.MsgRemovePriceFeedersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if ms.authority != msg.Authority {
		return nil, errors.Wrapf(
			types.ErrNoGovAuthority,
			"invalid authority; expected %s, got %s",
			ms.authority, msg.Authority,
		)
	}

	for _, feeder := range msg.Feeders {
		ms.Keeper.RemovePriceFeeder(ctx, sdk.MustAccAddressFromBech32(feeder))
	}
	return &types.MsgRemovePriceFeedersResponse{}, nil
}
