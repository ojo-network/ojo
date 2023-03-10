package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var Reasons = map[uint16]string{
	0: "No vote",
	1: "Vote above threshold",
	2: "Vote below threshold",
}

type (
	MissCounterLogs []MissCounterLog

	MissCounterLog struct {
		Operator sdk.ValAddress
		Reason   uint16
		Denom    string
	}
)

// ojoval1, ojoval2, ojoval3, and 20 others missed BTC/USDT vote
