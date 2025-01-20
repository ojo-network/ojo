package keeper

import (
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ojo-network/ojo/util"
	"github.com/ojo-network/ojo/x/oracle/types"
)

// SetPrice set a specific price in the store from its index
func (k Keeper) SetPrice(ctx sdk.Context, price types.Price) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := k.cdc.MustMarshal(&price)
	store.Set(types.KeyPrice(price.Asset, price.Source, price.Timestamp), bz)
}

// GetPrice returns a price from its index
func (k Keeper) GetPrice(ctx sdk.Context, asset, source string, timestamp uint64) (val types.Price, found bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	bz := store.Get(types.KeyPrice(asset, source, timestamp))
	if bz == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(bz, &val)
	return val, true
}

func (k Keeper) GetLatestPriceFromAssetAndSource(ctx sdk.Context, asset, source string) (val types.Price, found bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStoreReversePrefixIterator(store, types.KeyPriceAssetAndSource(asset, source))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Price
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		return val, true
	}

	return val, false
}

func (k Keeper) GetLatestPriceFromAnySource(ctx sdk.Context, asset string) (val types.Price, found bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStoreReversePrefixIterator(store, types.KeyPriceAsset(asset))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Price
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		return val, true
	}

	return val, false
}

// RemovePrice removes a price from the store
func (k Keeper) RemovePrice(ctx sdk.Context, asset, source string, timestamp uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(types.KeyPrice(asset, source, timestamp))
}

// GetAllPrice returns all price
func (k Keeper) GetAllPrice(ctx sdk.Context) (list []types.Price) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.KeyPrefixPrice)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Price
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) GetAssetPrice(ctx sdk.Context, asset string) (types.Price, bool) {
	return k.GetLatestPriceFromAnySource(ctx, asset)
}

func Pow10(decimal uint64) (value math.LegacyDec) {
	value = math.LegacyNewDec(1)
	for i := 0; i < util.SafeUint64ToInt(decimal); i++ {
		value = value.Mul(math.LegacyNewDec(10))
	}
	return
}

func (k Keeper) GetAssetPriceFromDenom(ctx sdk.Context, denom string) math.LegacyDec {
	info, found := k.GetAssetInfo(ctx, denom)
	if !found {
		return math.LegacyZeroDec()
	}
	price, found := k.GetAssetPrice(ctx, info.Display)
	if !found {
		return math.LegacyZeroDec()
	}
	return price.Price.Quo(Pow10(info.Decimal))
}

// SetAssetInfo set a specific assetInfo in the store from its index
func (k Keeper) SetAssetInfo(ctx sdk.Context, assetInfo types.AssetInfo) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := k.cdc.MustMarshal(&assetInfo)
	store.Set(types.KeyAssetInfo(assetInfo.Denom), bz)
}

// GetAssetInfo returns a assetInfo from its index
func (k Keeper) GetAssetInfo(ctx sdk.Context, denom string) (val types.AssetInfo, found bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := store.Get(types.KeyAssetInfo(denom))
	if bz == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(bz, &val)
	return val, true
}

// RemoveAssetInfo removes a assetInfo from the store
func (k Keeper) RemoveAssetInfo(ctx sdk.Context, denom string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(types.KeyAssetInfo(denom))
}

// GetAllAssetInfo returns all assetInfo
func (k Keeper) GetAllAssetInfo(ctx sdk.Context) (list []types.AssetInfo) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.KeyPrefixAssetInfo)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.AssetInfo
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// SetPriceFeeder set a specific priceFeeder in the store from its index
func (k Keeper) SetPriceFeeder(ctx sdk.Context, priceFeeder types.PriceFeeder) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.KeyPriceFeeder(priceFeeder.Feeder)
	bz := k.cdc.MustMarshal(&priceFeeder)
	store.Set(key, bz)
}

// GetPriceFeeder returns a priceFeeder from its index
func (k Keeper) GetPriceFeeder(ctx sdk.Context, feeder sdk.AccAddress) (val types.PriceFeeder, found bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.KeyPriceFeeder(feeder.String())

	bz := store.Get(key)
	if bz == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(bz, &val)
	return val, true
}

// RemovePriceFeeder removes a priceFeeder from the store
func (k Keeper) RemovePriceFeeder(ctx sdk.Context, feeder sdk.AccAddress) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.KeyPriceFeeder(feeder.String())
	store.Delete(key)
}

// SetPool sets a pool in the store from its index.
func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.KeyPool(pool.PoolId)
	bz := k.cdc.MustMarshal(&pool)
	store.Set(key, bz)
}

// GetPool returns a pool from its index
func (k Keeper) GetPool(ctx sdk.Context, poolID uint64) (val types.Pool, found bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.KeyPool(poolID)

	bz := store.Get(key)
	if bz == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(bz, &val)
	return val, true
}

// RemovePool removes a pool from the store
func (k Keeper) RemovePool(ctx sdk.Context, poolID uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.KeyPool(poolID)
	store.Delete(key)
}

// GetAllPool returns all pool
func (k Keeper) GetAllPool(ctx sdk.Context) (list []types.Pool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.KeyPrefixPool)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Pool
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// SetAccountedPool sets a accounted pool in the store from its index.
func (k Keeper) SetAccountedPool(ctx sdk.Context, accountedPool types.AccountedPool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.KeyAccountedPool(accountedPool.PoolId)
	bz := k.cdc.MustMarshal(&accountedPool)
	store.Set(key, bz)
}

// GetAccountedPool returns a accounted pool from its index
func (k Keeper) GetAccountedPool(ctx sdk.Context, poolID uint64) (val types.AccountedPool, found bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.KeyAccountedPool(poolID)

	bz := store.Get(key)
	if bz == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(bz, &val)
	return val, true
}

// RemoveAccountedPool removes a accounted pool from the store
func (k Keeper) RemoveAccountedPool(ctx sdk.Context, poolID uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.KeyAccountedPool(poolID)
	store.Delete(key)
}

// GetAllAccountedPool returns all accounted pool
func (k Keeper) GetAllAccountedPool(ctx sdk.Context) (list []types.AccountedPool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.KeyPrefixAccountedPool)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.AccountedPool
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
