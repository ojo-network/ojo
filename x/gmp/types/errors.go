package types

import "cosmossdk.io/errors"

var ErrNoPriceFound = errors.Register(ModuleName, 1, "no oracle price found")
