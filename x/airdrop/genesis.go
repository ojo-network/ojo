package airdrop

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/x/airdrop/keeper"
	"github.com/ojo-network/ojo/x/airdrop/types"
)

// InitGenesis initializes the x/airdrop module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, genState types.GenesisState) {
	err := keeper.SetParams(ctx, genState.Params)
	if err != nil {
		panic(err)
	}

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
	params, err := keeper.GetParams(ctx)
	if err != nil {
		panic(err)
	}
	genesisState.Params = params
	genesisState.AirdropAccounts = keeper.GetAllAirdropAccounts(ctx)
	return genesisState
}
