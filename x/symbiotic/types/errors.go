package types

import "cosmossdk.io/errors"

var (
	ErrSymbioticValUpdate = errors.Register(ModuleName, 1, "symbiotic validator update error")
	ErrSymbioticNotFound  = errors.Register(ModuleName, 2, "symbiotic not found")
)
