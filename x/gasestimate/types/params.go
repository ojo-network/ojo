package types

import (
	"fmt"
)

var (
	DefaultGasLimit         = "1000000"
	DefaultGasAdjustment    = "1.5"
	DefaultContractRegistry = []*Contract{
		{
			Address: "0x5BB3E85f91D08fe92a3D123EE35050b763D6E6A7",
			Network: "Ethereum",
		},
		{
			Address: "0x5BB3E85f91D08fe92a3D123EE35050b763D6E6A7",
			Network: "Arbitrum",
		},
	}
)

func DefaultParams() Params {
	return Params{
		ContractRegistry: DefaultContractRegistry,
		GasLimit:         DefaultGasLimit,
		GasAdjustment:    DefaultGasAdjustment,
	}
}

func (p Params) Validate() error {
	if len(p.ContractRegistry) == 0 {
		return fmt.Errorf("contract registry can not be empty")
	}
	if p.GasLimit == "" {
		return fmt.Errorf("gas limit can not be empty")
	}
	if p.GasAdjustment == "" {
		return fmt.Errorf("gas adjustment can not be empty")
	}
	return nil
}
