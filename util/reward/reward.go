package reward

import (
	"fmt"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CalculateRewardFactor returns the reward factor calculated using a logarmithic
// model based on miss counters. missCount is the current miss count, m is the
// maximum possible miss counts, and s is the smallest miss count in the period.
// If the logarimthic function returns NaN the Reward Factor returned will be 0.
// rewardFactor = 1 - logₘ₋ₛ₊₁(missCount - s + 1)
func CalculateRewardFactor(missCount sdk.Dec, m sdk.Dec, s sdk.Dec) string {
	logBase := m.Sub(s).Add(sdk.NewDec(1))
	logKey := missCount.Sub(s).Add(sdk.NewDec(1))
	rewardFactor := 1 - (math.Log(logKey.MustFloat64()) / math.Log(logBase.MustFloat64()))
	if math.IsNaN(rewardFactor) {
		rewardFactor = 0
	}

	return fmt.Sprintf("%f", rewardFactor)
}
