package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ojo-network/ojo/x/airdrop/types"
)

// SetParams sets the airdrop module's parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) (err error) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ParamsKey, k.cdc.MustMarshal(&params))
	return
}

// GetParams gets the airdrop module's parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params, err error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return types.DefaultParams(), nil
	}
	k.cdc.MustUnmarshal(bz, &params)
	return
}
