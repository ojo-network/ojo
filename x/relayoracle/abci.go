package relayoracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/x/relayoracle/keeper"
)

// Process all requests in the endblocker, after oracle module end blocker
// EndBlocker resolves all the pending requests
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// Loops through all requests in the resolvable list to resolve all of them!
	for _, reqID := range k.GetPendingRequestList(ctx) {
		k.ResolveRequest(ctx, reqID)
	}

	// Once all the requests are resolved, we can clear the list.
	k.FlushPendingRequestList(ctx)
}
