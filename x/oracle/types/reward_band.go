package types

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// String implements fmt.Stringer interface
func (rb RewardBand) String() string {
	out, _ := yaml.Marshal(rb)
	return string(out)
}

func (rb RewardBand) Equal(rb2 *RewardBand) bool {
	if !strings.EqualFold(rb.SymbolDenom, rb2.SymbolDenom) {
		return false
	}
	if !rb.RewardBand.Equal(rb2.RewardBand) {
		return false
	}
	return true
}

// RewardBandList is array of RewardBand
type RewardBandList []RewardBand

func (rbl RewardBandList) String() (out string) {
	for _, d := range rbl {
		out += d.String() + "\n"
	}

	return strings.TrimSpace(out)
}
