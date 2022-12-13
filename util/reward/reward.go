package reward

import (
	"fmt"
	"math"
)

func CalculateRewardFactor(missCount uint64, m uint64, s uint64) string {
	logBase := m - s + 1
	rewardFactor := 1 - (math.Log(float64(missCount-s+1)) / math.Log(float64(logBase)))

	return fmt.Sprintf("%f", rewardFactor)
}
