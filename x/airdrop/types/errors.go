package types

import "cosmossdk.io/errors"

var (
	ErrAirdropAlreadyClaimed  = errors.Register(ModuleName, 2, "airdrop account already claimed")
	ErrInsufficientDelegation = errors.Register(ModuleName, 3, "delegation requirement not met")
	ErrAirdropExpired         = errors.Register(ModuleName, 4, "airdrop expired; chain is past the expire block")
)
