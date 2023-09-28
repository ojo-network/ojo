package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/relayoracle module sentinel errors
var (
	ErrInvalidPacketTimeout = sdkerrors.Register(ModuleName, 1, "invalid packet timeout")
	ErrInvalidVersion       = sdkerrors.Register(ModuleName, 2, "invalid version")
	ErrIbcRequestDisabled   = sdkerrors.Register(ModuleName, 3, "ibc request.go disabled")
	ErrRequestNotFound      = sdkerrors.Register(ModuleName, 4, "request not found")
	ErrTooManyDenoms        = sdkerrors.Register(ModuleName, 5, "total denoms exceeds threshold")
	ErrNoDenoms             = sdkerrors.Register(ModuleName, 5, "no denoms in reuqest")
)
