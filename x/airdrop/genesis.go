package airdrop

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/x/airdrop/keeper"
	"github.com/ojo-network/ojo/x/airdrop/types"
)

// InitGenesis initializes the x/airdrop module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, genState types.GenesisState) {
	keeper.SetParams(ctx, genState.Params)

	for _, airdropAccount := range genState.AirdropAccounts {
		err := keeper.SetAirdropAccount(ctx, airdropAccount)
		if err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns the x/airdrop module's exported genesis.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	genesisState := types.DefaultGenesisState()
	genesisState.Params = keeper.GetParams(ctx)
	genesisState.AirdropAccounts = keeper.GetAllAirdropAccounts(ctx)
	return genesisState
}
