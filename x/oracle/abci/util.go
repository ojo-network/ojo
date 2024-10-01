package abci

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

// VoteExtensionsEnabled determines if vote extensions are enabled for the current block.
func VoteExtensionsEnabled(ctx sdk.Context) bool {
	cp := ctx.ConsensusParams()
	if cp.Abci == nil || cp.Abci.VoteExtensionsEnableHeight == 0 {
		return false
	}

	// Per the cosmos sdk, the first block should not utilize the latest finalize block state. This means
	// vote extensions should NOT be making state changes.
	//
	// Ref: https://github.com/cosmos/cosmos-sdk/blob/2100a73dcea634ce914977dbddb4991a020ee345/baseapp/baseapp.go#L488-L495
	if ctx.BlockHeight() <= 1 {
		return false
	}

	return cp.Abci.VoteExtensionsEnableHeight < ctx.BlockHeight()
}

func calculateMedian(gasEstimates []oracletypes.GasEstimate) (median oracletypes.GasEstimate) {
	sort.Slice(gasEstimates, func(i, j int) bool {
		return gasEstimates[i].GasEstimation < gasEstimates[j].GasEstimation
	})

	mid := len(gasEstimates) / 2
	if len(gasEstimates)%2 == 0 {
		return oracletypes.GasEstimate{
			GasEstimation: gasEstimates[mid-1].GasEstimation + gasEstimates[mid].GasEstimation,
			Network:       gasEstimates[mid-1].Network,
		}
	}
	return gasEstimates[mid]
}
