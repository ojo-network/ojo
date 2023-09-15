package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/relayoracle module sentinel errors
var (
	ErrSample               = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrInvalidPacketTimeout = sdkerrors.Register(ModuleName, 1500, "invalid packet timeout")
	ErrInvalidVersion       = sdkerrors.Register(ModuleName, 1501, "invalid version")
	ErrIbcRequestDisabled   = sdkerrors.Register(ModuleName, 1502, "ibc request.go disabled")
	ErrRequestNotFound      = sdkerrors.Register(ModuleName, 1503, "request not found")
	ErrNoActiveExchangeRate = sdkerrors.Register(ModuleName, 1504, "no active exchange rate")
)
