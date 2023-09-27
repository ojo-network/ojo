package relayoracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/x/relayoracle/keeper"
)

// Process all requests in the endblocker, after oracle module end blocker
// EndBlocker resolves all the pending requests
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// Resolve all the pending requests
	for _, reqID := range k.GetPendingRequestList(ctx) {
		k.ResolveRequest(ctx, reqID)
	}

	// Clear all pending requests
	k.FlushPendingRequestList(ctx)
}
