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
}

// ExportGenesis returns the module's exported genesis
// TODO: add pending request list and results to genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	genesis.PortId = k.GetPort(ctx)

	return genesis
}
