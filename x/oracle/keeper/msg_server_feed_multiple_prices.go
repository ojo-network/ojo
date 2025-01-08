package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ojo-network/ojo/x/oracle/types"
)

func (ms msgServer) FeedMultiplePrices(goCtx context.Context, msg *types.MsgFeedMultiplePrices) (*types.MsgFeedMultiplePricesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	feeder, found := ms.Keeper.GetPriceFeeder(ctx, creator)
	if !found {
		return nil, types.ErrNotAPriceFeeder
	}

	if !feeder.IsActive {
		return nil, types.ErrPriceFeederNotActive
	}

	for _, feedPrice := range msg.FeedPrices {
		price := types.Price{
			Asset:       feedPrice.Asset,
			Price:       feedPrice.Price,
			Source:      feedPrice.Source,
			Provider:    msg.Creator,
			Timestamp:   uint64(ctx.BlockTime().Unix()),
			BlockHeight: uint64(ctx.BlockHeight()),
		}
		ms.SetPrice(ctx, price)
	}

	return &types.MsgFeedMultiplePricesResponse{}, nil
}
