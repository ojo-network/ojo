package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var Reasons = map[uint16]string{
	0: "No vote",
	1: "Vote above threshold",
	2: "Vote below threshold",
}

type (
	// MissCounterLog
	MissCounterLogs map[string]MissCounterLog

	MissCounterLog struct {
		Violators   []sdk.ValAddress
		Reason      uint16
		Denom       string
		VotePrice   sdk.Dec
		MedianPrice sdk.Dec
	}
)

func (mcl MissCounterLog) Key() string {
	return fmt.Sprintf("%d:%s", mcl.Reason, mcl.Denom)
}

// OutputString returns a string representation of the miss counter log
// where the first three violators are listed, total number of
// violators, along with the reason and denom.
func (mcl MissCounterLog) OutputString() string {
	violators := ""
	for i, violator := range mcl.Violators {
		if i == 3 {
			break
		}
		violators += violator.String() + ", "
	}
	if len(mcl.Violators) > 3 {
		violators += "..."
	}
	return fmt.Sprintf("%s %s %s %s", violators, mcl.VotePrice, mcl.MedianPrice, Reasons[mcl.Reason])
}
