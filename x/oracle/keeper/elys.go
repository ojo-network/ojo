package keeper

import (
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ojo-network/ojo/x/oracle/types"
)

// SetPrice set a specific price in the store from its index
func (k Keeper) SetPrice(ctx sdk.Context, price types.Price) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&price)
	store.Set(types.KeyPrice(price.Asset, price.Source, price.Timestamp), bz)
}

// GetPrice returns a price from its index
func (k Keeper) GetPrice(ctx sdk.Context, asset, source string, timestamp uint64) (val types.Price, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyPrice(asset, source, timestamp))
	if bz == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(bz, &val)
	return val, true
}

func (k Keeper) GetLatestPriceFromAssetAndSource(ctx sdk.Context, asset, source string) (val types.Price, found bool) {
	store := ctx.KVStore(k.storeKey)
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
	store := ctx.KVStore(k.storeKey)
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
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyPrice(asset, source, timestamp))
}

// GetAllPrice returns all price
func (k Keeper) GetAllPrice(ctx sdk.Context) (list []types.Price) {
	store := ctx.KVStore(k.storeKey)
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
	for i := 0; i < int(decimal); i++ {
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
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&assetInfo)
	store.Set(types.KeyAssetInfo(assetInfo.Denom), bz)
}

// GetAssetInfo returns a assetInfo from its index
func (k Keeper) GetAssetInfo(ctx sdk.Context, denom string) (val types.AssetInfo, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyAssetInfo(denom))
	if bz == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(bz, &val)
	return val, true
}

// RemoveAssetInfo removes a assetInfo from the store
func (k Keeper) RemoveAssetInfo(ctx sdk.Context, denom string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyAssetInfo(denom))
}

// GetAllAssetInfo returns all assetInfo
func (k Keeper) GetAllAssetInfo(ctx sdk.Context) (list []types.AssetInfo) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

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
	store := ctx.KVStore(k.storeKey)
	key := types.KeyPriceFeeder(priceFeeder.Feeder)
	bz := k.cdc.MustMarshal(&priceFeeder)
	store.Set(key, bz)
}

// GetPriceFeeder returns a priceFeeder from its index
func (k Keeper) GetPriceFeeder(ctx sdk.Context, feeder sdk.AccAddress) (val types.PriceFeeder, found bool) {
	store := ctx.KVStore(k.storeKey)
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
	store := ctx.KVStore(k.storeKey)
	key := types.KeyPriceFeeder(feeder.String())
	store.Delete(key)
}
