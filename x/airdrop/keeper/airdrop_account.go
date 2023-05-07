package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ojo-network/ojo/x/airdrop/types"
)

// SetAirdropAccount saves the airdrop account to the store
// using the OriginAddress as the key.
func (k Keeper) SetAirdropAccount(
	ctx sdk.Context,
	account types.AirdropAccount,
) (err error) {
	store := ctx.KVStore(k.storeKey)
	store.Set(
		types.AirdropAccountKey(account.OriginAddress),
		k.cdc.MustMarshal(&account),
	)
	return
}

// GetAirdropAccount returns the airdrop account from the store
func (k Keeper) GetAirdropAccount(
	ctx sdk.Context, originAddress string,
) (account types.AirdropAccount, err error) {

	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AirdropAccountKey(originAddress))
	if bz == nil {
		return
	}

	k.cdc.MustUnmarshal(bz, &account)
	return
}

// GetAllAirdropAccounts returns all airdrop accounts from the store
func (k Keeper) GetAllAirdropAccounts(
	ctx sdk.Context,
) (accounts []types.AirdropAccount) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.AirdropAccountKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var account types.AirdropAccount
		k.cdc.MustUnmarshal(iterator.Value(), &account)
		accounts = append(accounts, account)
	}
	return
}
