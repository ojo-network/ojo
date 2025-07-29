package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ojo-network/ojo/x/oracle/types"
)

func (ms msgServer) CreateAssetInfo(
	goCtx context.Context,
	msg *types.MsgCreateAssetInfo,
) (*types.MsgCreateAssetInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, found := ms.GetAssetInfo(ctx, msg.Denom)

	if found {
		return nil, errors.Wrapf(types.ErrAssetWasCreated, "%s", msg.Denom)
	}

	ms.Keeper.SetAssetInfo(ctx, types.AssetInfo{
		Denom:   msg.Denom,
		Display: msg.Display,
		Ticker:  msg.Ticker,
		Decimal: msg.Decimal,
	})

	return &types.MsgCreateAssetInfoResponse{}, nil
}

func (ms msgServer) RemoveAssetInfo(
	goCtx context.Context,
	msg *types.MsgRemoveAssetInfo,
) (*types.MsgRemoveAssetInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if ms.authority != msg.Authority {
		return nil, errors.Wrapf(
			types.ErrNoGovAuthority,
			"invalid authority; expected %s, got %s",
			ms.authority, msg.Authority,
		)
	}

	ms.Keeper.RemoveAssetInfo(ctx, msg.Denom)
	return &types.MsgRemoveAssetInfoResponse{}, nil
}
