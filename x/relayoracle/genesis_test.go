package relayoracle_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/ojo-network/ojo/testutil/keeper"
	"github.com/ojo-network/ojo/x/relayoracle"
	"github.com/ojo-network/ojo/x/relayoracle/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
	}

	k, ctx := keepertest.RelayoracleKeeper(t)
	relayoracle.InitGenesis(ctx, *k, genesisState)
	got := relayoracle.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	require.Equal(t, genesisState.PortId, got.PortId)
}
