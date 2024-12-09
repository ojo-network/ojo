package keeper

import (
	"fmt"

	"cosmossdk.io/log"

	sdk "github.com/cosmos/cosmos-sdk/types"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ojo-network/ojo/x/symbiotic/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
	apiUrls  types.ApiUrls

	StakingKeeper types.StakingKeeper

	// the address capable of executing a MsgSetParams message. Typically, this
	// should be the x/gov module account.
	authority string
}

// NewKeeper constructs a new keeper for gmp module.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	authority string,
) Keeper {
	return Keeper{
		cdc:       cdc,
		storeKey:  storeKey,
		apiUrls:   types.NewApiUrls(),
		authority: authority,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) SetCachedBlockHash(
	ctx sdk.Context,
	cachedBlockHash types.CachedBlockHash,
) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&cachedBlockHash)
	store.Set(types.CachedBlockHashKey(uint64(cachedBlockHash.Height)), bz)
}

func (k Keeper) DeleteCachedBlockHash(
	ctx sdk.Context,
	cachedBlockHash types.CachedBlockHash,
) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.CachedBlockHashKey(uint64(cachedBlockHash.Height)))
}

func (k Keeper) GetCachedBlockHash(
	ctx sdk.Context,
	blockHeight uint64,
) (types.CachedBlockHash, error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.CachedBlockHashKey(blockHeight))
	cachedBlockHash := types.CachedBlockHash{}
	k.cdc.MustUnmarshal(bz, &cachedBlockHash)
	return cachedBlockHash, nil
}

func (k Keeper) IterateAllCachedBlockHashes(
	ctx sdk.Context,
	handler func(types.CachedBlockHash) bool,
) {
	store := ctx.KVStore(k.storeKey)
	iter := storetypes.KVStorePrefixIterator(store, types.CachedBlockHashPrefix)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		cachedBlockHash := types.CachedBlockHash{}
		k.cdc.MustUnmarshal(iter.Value(), &cachedBlockHash)
		if handler(cachedBlockHash) {
			break
		}
	}
}

func (k Keeper) GetAllCachedBlockHashes(
	ctx sdk.Context,
) []types.CachedBlockHash {
	cachedBlockHashes := []types.CachedBlockHash{}
	k.IterateAllCachedBlockHashes(ctx, func(cachedBlockHash types.CachedBlockHash) (stop bool) {
		cachedBlockHashes = append(cachedBlockHashes, cachedBlockHash)
		return false
	})
	return cachedBlockHashes
}

func (k Keeper) PruneBlockHashesBeforeBlock(
	ctx sdk.Context,
	blockNum uint64,
) {
	k.IterateAllCachedBlockHashes(ctx, func(cachedBlockHash types.CachedBlockHash) (stop bool) {
		if cachedBlockHash.Height <= int64(blockNum) {
			k.DeleteCachedBlockHash(ctx, cachedBlockHash)
		}
		return false
	})
}
