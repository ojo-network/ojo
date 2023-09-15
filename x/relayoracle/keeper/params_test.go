package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	testkeeper "github.com/ojo-network/ojo/testutil/keeper"
	"github.com/ojo-network/ojo/x/relayoracle/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.RelayoracleKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
