package types

import "cosmossdk.io/errors"

var (
	ErrNoAccountFound         = errors.Register(ModuleName, 2, "no airdrop account found")
	ErrAirdropAlreadyClaimed  = errors.Register(ModuleName, 3, "airdrop account already claimed")
	ErrInsufficientDelegation = errors.Register(ModuleName, 4, "delegation requirement not met")
	ErrAirdropExpired         = errors.Register(ModuleName, 5, "airdrop expired; chain is past the expire block")
	ErrOriginAccountExists    = errors.Register(ModuleName, 6, "origin account already exists")
)
