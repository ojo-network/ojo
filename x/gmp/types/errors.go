package types

import "cosmossdk.io/errors"

var ErrInvalidDestinationChain = errors.Register(ModuleName, 1, "invalid destination chain")
