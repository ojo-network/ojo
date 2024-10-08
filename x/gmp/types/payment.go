package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p Payment) TriggerUpdate(rate math.LegacyDec, ctx sdk.Context) bool {
	// Calculate the percentage difference
	priceDiff := p.LastPrice.Sub(rate).Abs()
	percentageDiff := priceDiff.Quo(p.LastPrice).MulInt64(100)

	return percentageDiff.GT(p.Deviation) ||
		p.LastBlock < ctx.BlockHeight()-p.Heartbeat
}
