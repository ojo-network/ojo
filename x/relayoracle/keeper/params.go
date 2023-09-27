package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/x/relayoracle/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams()
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

func (k Keeper) SetIbcRequestEnabled(ctx sdk.Context, enabled bool) {
	k.paramstore.Set(ctx, types.KeyIbcRequestEnabled, enabled)
}

func (k Keeper) SetPacketExpiryBlockCount(ctx sdk.Context, expiry uint64) {
	k.paramstore.Set(ctx, types.KeyPacketExpiry, expiry)
}

func (k Keeper) IbcRequestEnabled(ctx sdk.Context) (enabled bool) {
	k.paramstore.Get(ctx, types.KeyIbcRequestEnabled, &enabled)
	return
}

func (k Keeper) PacketExpiry(ctx sdk.Context) (expiry uint64) {
	k.paramstore.Get(ctx, types.KeyPacketExpiry, &expiry)
	return
}
