package types

import "cosmossdk.io/errors"

var (
	ErrInvalidDestinationChain = errors.Register(ModuleName, 1, "invalid destination chain")
	ErrEncodeInjVoteExt        = errors.Register(ModuleName, 2, "failed to encode injected vote extension tx")
	ErrNoCommitInfo            = errors.Register(ModuleName, 3, "no commit info")
)
