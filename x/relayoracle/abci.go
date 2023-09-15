package relayoracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/x/relayoracle/keeper"
)

// TODO: this might not be necessary, could return the data immediately?
// bcz, the new rates for oracle updated at their end block
// // handleEndBlock cleans up the state during end block. See comment in the implementation!
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// Loops through all requests in the resolvable list to resolve all of them!
	for _, reqID := range k.GetPendingRequestList(ctx) {
		k.ResolveRequest(ctx, reqID)
	}

	// Once all the requests are resolved, we can clear the list.
	k.FlushPendingRequestList(ctx)
}
