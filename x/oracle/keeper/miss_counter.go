package keeper

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/ojo-network/ojo/x/oracle/types"
)

// GetMissCounter retrieves the # of vote periods missed in this oracle slash
// window.
func (k Keeper) GetMissCounter(ctx sdk.Context, operator string) uint64 {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	bz := store.Get(types.GetMissCounterKey(operator))
	if bz == nil {
		// by default the counter is zero
		return 0
	}

	var missCounter gogotypes.UInt64Value
	k.cdc.MustUnmarshal(bz, &missCounter)

	return missCounter.Value
}

// SetMissCounter updates the # of vote periods missed in this oracle slash
// window.
func (k Keeper) SetMissCounter(ctx sdk.Context, operator string, missCounter uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: missCounter})
	store.Set(types.GetMissCounterKey(operator), bz)
}

// DeleteMissCounter removes miss counter for the validator.
func (k Keeper) DeleteMissCounter(ctx sdk.Context, operator string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(types.GetMissCounterKey(operator))
}

// IterateMissCounters iterates over the miss counters and performs a callback
// function.
func (k Keeper) IterateMissCounters(ctx sdk.Context, handler func(string, uint64) bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	iter := storetypes.KVStorePrefixIterator(store, types.KeyPrefixMissCounter)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		operator := string(iter.Key()[1:])

		var missCounter gogotypes.UInt64Value
		k.cdc.MustUnmarshal(iter.Value(), &missCounter)

		if handler(operator, missCounter.Value) {
			break
		}
	}
}
