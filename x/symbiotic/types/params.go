package types

import (
	"fmt"
)

var (
	DefaultMiddlewareAddress = "0x0000000000000000000000000000000000000000"
)

func DefaultParams() Params {
	return Params{
		MiddlewareAddress: DefaultMiddlewareAddress,
	}
}

func (p Params) Validate() error {
	if len(p.MiddlewareAddress) != 42 {
		return fmt.Errorf("middleware address must be a valid ethereum address")
	}
	return nil
}
