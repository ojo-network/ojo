package gas_estimate

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/x/gas_estimate/keeper"
	"github.com/ojo-network/ojo/x/gas_estimate/types"
)

// InitGenesis initializes the x/gas_estimate module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, genState types.GenesisState) {
	keeper.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the x/gas_estimate module's exported genesis.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	genesisState := types.DefaultGenesisState()
	genesisState.Params = keeper.GetParams(ctx)
	return genesisState
}
