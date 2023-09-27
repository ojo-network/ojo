package types

import (
	"fmt"
	"time"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

const (
	DefaultIbcRequestEnabled = true
	DefaultPacketTimeout     = uint64(10 * time.Minute)
)

var (
	KeyIbcRequestEnabled = []byte("IbcRequestEnabled")
	KeyPacketExpiry      = []byte("PacketExpiry")
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		IbcRequestEnabled: DefaultIbcRequestEnabled,
		PacketTimeout:     DefaultPacketTimeout,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyIbcRequestEnabled, &p.IbcRequestEnabled, validateBool),
		paramtypes.NewParamSetPair(KeyPacketExpiry, &p.PacketTimeout, validateUint64),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	return nil
}

func validateBool(value interface{}) error {
	_, ok := value.(bool)
	if !ok {
		return fmt.Errorf("invalid  value: %T", value)
	}

	return nil
}

func validateUint64(value interface{}) error {
	v, ok := value.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", value)
	}
	if v <= 0 {
		return fmt.Errorf("value must be positive: %T", value)
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
