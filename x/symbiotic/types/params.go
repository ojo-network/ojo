package types

import (
	"fmt"
)

var (
	DefaultMiddlewareAddress        = "0x0000000000000000000000000000000000000000"
	DefaultSymbioticSyncPeriod      = int64(10)
	DefaultMaximumCachedBlockHashes = uint64(10)
)

func DefaultParams() Params {
	return Params{
		MiddlewareAddress:        DefaultMiddlewareAddress,
		SymbioticSyncPeriod:      DefaultSymbioticSyncPeriod,
		MaximumCachedBlockHashes: DefaultMaximumCachedBlockHashes,
	}
}

func (p Params) Validate() error {
	if len(p.MiddlewareAddress) != 42 {
		return fmt.Errorf("middleware address must be a valid ethereum address")
	}
	if p.SymbioticSyncPeriod < 1 {
		return fmt.Errorf("symbiotic sync period value must be positive and nonzero")
	}
	return nil
}
