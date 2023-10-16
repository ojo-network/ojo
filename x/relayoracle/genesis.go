package relayoracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/x/relayoracle/keeper"
	"github.com/ojo-network/ojo/x/relayoracle/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
// TODO: export pending list and other state to genesis
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetPort(ctx, genState.PortId)
	if !k.IsBound(ctx, genState.PortId) {
		err := k.BindPort(ctx, genState.PortId)
		if err != nil {
			panic("could not claim port capability: " + err.Error())
		}
	}
	k.SetRequestCount(ctx, 0)

	k.SetParams(ctx, genState.Params)

	// set request ids
	for _, request := range genState.Requests {
		k.SetRequest(ctx, request.RequestID, request)
	}

	// set results
	for _, result := range genState.Results {
		k.SetResult(ctx, result)
	}

	// set pending list
	var pending types.PendingRequestList
	pending.RequestIds = genState.GetPendingRequestIds()
	k.SetPendingList(ctx, pending)

}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	genesis.PortId = k.GetPort(ctx)

	// extract all results, requests and pending request ids
	genesis.Results = k.AllResults(ctx)
	genesis.Requests = k.AllRequests(ctx)
	genesis.PendingRequestIds = k.GetPendingRequestList(ctx)

	return genesis
}
