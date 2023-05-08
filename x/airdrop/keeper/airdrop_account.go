package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"

	"github.com/ojo-network/ojo/x/airdrop/types"
)

// SetAirdropAccount saves the airdrop account to the store
// using the OriginAddress as the key.
func (k Keeper) SetAirdropAccount(
	ctx sdk.Context,
	account *types.AirdropAccount,
) (err error) {
	store := ctx.KVStore(k.storeKey)
	store.Set(
		types.AirdropAccountKey(account.OriginAddress),
		k.cdc.MustMarshal(account),
	)
	return
}

// GetAirdropAccount returns the airdrop account from the store
func (k Keeper) GetAirdropAccount(
	ctx sdk.Context, originAddress string,
) (account *types.AirdropAccount, err error) {

	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AirdropAccountKey(originAddress))
	if bz == nil {
		return
	}

	k.cdc.MustUnmarshal(bz, account)
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

// VerifyDelegationRequirement returns an error if the total shares
// delegated is less than the delegation requirement.
func (k Keeper) VerifyDelegationRequirement(
	ctx sdk.Context,
	aa *types.AirdropAccount,
) (err error) {
	delegations := k.stakingKeeper.GetDelegatorDelegations(ctx, aa.ClaimAccAddress(), 999)
	totalShares := sdk.ZeroDec()
	for _, delegation := range delegations {
		totalShares = totalShares.Add(delegation.Shares)
	}
	if totalShares.LT(*k.GetParams(ctx).DelegationRequirement) {
		return types.ErrInsufficientDelegation
	}
	return nil
}

// SetClaimAmount calculates and sets the claim amount for the airdrop account
func (k Keeper) SetClaimAmount(ctx sdk.Context, aa *types.AirdropAccount) {
	claimAmount := k.GetParams(ctx).AirdropFactor.MulInt64(int64(aa.OriginAmount))
	aa.ClaimAmount = claimAmount.TruncateInt().Uint64()
}

// MintOriginTokens mints the originAmount of tokens to the airdrop module account
func (k Keeper) MintOriginTokens(ctx sdk.Context, aa *types.AirdropAccount) {
	k.bankKeeper.MintCoins(ctx, types.ModuleName, aa.OriginCoins())
}

// MintClaimTokens mints the claimAmount of tokens to the airdrop module account
func (k Keeper) MintClaimTokens(ctx sdk.Context, aa *types.AirdropAccount) {
	k.bankKeeper.MintCoins(ctx, types.ModuleName, aa.ClaimCoins())
}

func (k Keeper) CreateOriginAccount(ctx sdk.Context, aa *types.AirdropAccount) {
	baseAccount := authtypes.NewBaseAccountWithAddress(aa.OriginAccAddress())
	baseAccount = k.accountKeeper.NewAccount(ctx, baseAccount).(*authtypes.BaseAccount)
	baseVestingAccount := authvesting.NewBaseVestingAccount(baseAccount, aa.OriginCoins().Sort(), aa.VestingEndTime)
	vestingAccount := authvesting.NewContinuousVestingAccountRaw(baseVestingAccount, ctx.BlockTime().Unix())
	k.accountKeeper.SetAccount(ctx, vestingAccount)
}

func (k Keeper) CreateClaimAccount(ctx sdk.Context, aa *types.AirdropAccount) {
	baseAccount := authtypes.NewBaseAccountWithAddress(aa.ClaimAccAddress())
	baseAccount = k.accountKeeper.NewAccount(ctx, baseAccount).(*authtypes.BaseAccount)
	baseVestingAccount := authvesting.NewBaseVestingAccount(baseAccount, aa.ClaimCoins().Sort(), aa.VestingEndTime)
	vestingAccount := authvesting.NewDelayedVestingAccountRaw(baseVestingAccount)
	k.accountKeeper.SetAccount(ctx, vestingAccount)
}

func (k Keeper) SendOriginTokens(ctx sdk.Context, aa *types.AirdropAccount) {
	k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, aa.OriginAccAddress(), aa.OriginCoins())
}

func (k Keeper) SendClaimTokens(ctx sdk.Context, aa *types.AirdropAccount) {
	k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, aa.ClaimAccAddress(), aa.ClaimCoins())
}
