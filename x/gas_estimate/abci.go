package gas_estimate

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/x/gas_estimate/keeper"
)

// EndBlocker is called at the end of every block
func EndBlocker(_ sdk.Context, _ keeper.Keeper) error {
	return nil
}
