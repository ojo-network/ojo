package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p Payment) TriggerUpdate(rate math.LegacyDec, ctx sdk.Context) bool {
	return p.LastPrice.Sub(rate).Abs().GT(p.Deviation) ||
		p.LastBlock < ctx.BlockHeight()-p.Heartbeat
}
