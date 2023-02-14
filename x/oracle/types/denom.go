package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"gopkg.in/yaml.v3"
)

// String implements fmt.Stringer interface
func (d Denom) String() string {
	out, _ := yaml.Marshal(d)
	return string(out)
}

// Equal implements equal interface
func (d Denom) Equal(d1 *Denom) bool {
	return d.BaseDenom == d1.BaseDenom &&
		d.SymbolDenom == d1.SymbolDenom &&
		d.Exponent == d1.Exponent
}

// DenomList is array of Denom
type DenomList []Denom

// String implements fmt.Stringer interface
func (dl DenomList) String() (out string) {
	for _, d := range dl {
		out += d.String() + "\n"
	}

	return strings.TrimSpace(out)
}

// Contains checks whether or not a SymbolDenom (e.g. OJO) is in the DenomList
func (dl DenomList) Contains(symbolDenom string) bool {
	for _, d := range dl {
		if strings.EqualFold(d.SymbolDenom, symbolDenom) {
			return true
		}
	}
	return false
}

// ContainDenoms checks if d is a subset of dl
func (dl DenomList) ContainDenoms(d DenomList) bool {
	contains := make(map[string]struct{})

	for _, denom := range dl {
		contains[denom.String()] = struct{}{}
	}

	for _, denom := range d {
		if _, found := contains[denom.String()]; !found {
			return false
		}
	}

	return true
}

// GetRewardBand returns the reward band of a given asset in the DenomList.
// It will return an error if it can not find it.
func (dl DenomList) GetRewardBand(rbl RewardBandList) (sdk.Dec, error) {
	for _, d := range dl {
		for _, rb := range rbl {
			if strings.ToUpper(d.SymbolDenom) == strings.ToUpper(rb.SymbolDenom) {
				return rb.RewardBand, nil
			}
		}
	}
	return sdk.ZeroDec(), ErrNoRewardBand
}
