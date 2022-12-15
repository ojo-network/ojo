package reward

import (
	"fmt"
	"math"
)

// CalculateRewardFactor returns the reward factor calculated using a logarmithic
// model based on miss counters. missCount is the current miss count, m is the
// maximum possible miss counts, and s is the smallest miss count in the period.
// rewardFactor = 1 - logₘ₋ₛ₊₁(missCount - s + 1)
func CalculateRewardFactor(missCount uint64, m uint64, s uint64) string {
	logBase := m - s + 1
	rewardFactor := 1 - (math.Log(float64(missCount-s+1)) / math.Log(float64(logBase)))

	return fmt.Sprintf("%f", rewardFactor)
}