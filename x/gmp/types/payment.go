package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TriggerUpdate checks if the payment should be updated based on the current rate and the last update block.
func (p Payment) TriggerUpdate(rate math.LegacyDec, ctx sdk.Context) bool {
	if p.LastBlock == 0 || p.LastPrice.IsZero() {
		return true
	}
	priceDiff := p.LastPrice.Sub(rate).Abs()
	deviationExceeded := priceDiff.Quo(p.LastPrice).MulInt64(100).GT(p.Deviation)

	return deviationExceeded || p.LastBlock < ctx.BlockHeight()-p.Heartbeat
}
