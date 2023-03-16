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
		Violators []sdk.ValAddress
		Reason    uint16
		Denom     string
	}
)

func (mcls *MissCounterLogs) Add(address sdk.ValAddress, reason uint16, denom string) {
	logs := *mcls
	key := fmt.Sprintf("%d:%s", reason, denom)
	if _, ok := logs[key]; !ok {
		logs[key] = MissCounterLog{
			Violators: []sdk.ValAddress{address},
			Reason:    reason,
			Denom:     denom,
		}
		return
	} else {
		logs[key].Violators = append(logs[key].Violators, address)
	}
}

func (mcl MissCounterLog) Key() string {
	return fmt.Sprintf("%d:%s", mcl.Reason, mcl.Denom)
}

// String returns a string representation of the miss counter log
// where the first three violators are listed, total number of
// violators, along with the reason and denom.
func (mcl MissCounterLog) String() string {
	violators := ""
	for i, violator := range mcl.Violators {
		if i == 3 {
			break
		}
		violators += violator.String() + ", "
	}
	if len(mcl.Violators) > 3 {
		violators += fmt.Sprintf("...%d more", len(mcl.Violators))
	}
	return fmt.Sprintf("%s %s %s", violators, mcl.Denom, Reasons[mcl.Reason])
}
