package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ojo-network/ojo/x/gasestimate/types"
)

// SetParams sets the gasestimate module's parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ParamsKey, k.cdc.MustMarshal(&params))
}

// GetParams gets the gasestimate module's parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	k.cdc.MustUnmarshal(bz, &params)
	return
}
