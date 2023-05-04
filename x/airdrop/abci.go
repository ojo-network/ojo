package airdrop

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/x/airdrop/keeper"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) error {
	// TODO Check for ExpiryBlock
	// Query the number of unclaimed accounts?
	// all unclaimed AirdropAccounts will instead mint tokens into the community pool.
	return nil
}
