package util

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// IsPeriodLastBlock returns true if we are at the last block of the period
func IsPeriodLastBlock(ctx sdk.Context, blocksPerPeriod uint64) bool {
	return (SafeInt64ToUint64(ctx.BlockHeight())+1)%blocksPerPeriod == 0
}
