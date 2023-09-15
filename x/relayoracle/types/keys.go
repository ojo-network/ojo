package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "relayoracle"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_relayoracle"

	// Version defines the current version the IBC module supports
	Version = "relayoracle-1"

	// PortID is the default port id that module binds to
	PortID = "relayoracle"
)

var (
	// PortKey defines the key to store the port ID in store
	PortKey      = KeyPrefix("relayoracle-port-")
	RequestIDKey = KeyPrefix("relayoracle-request.go-")

	RequestStoreKeyPrefix = []byte{0x01}
	ResultStoreKeyPrefix  = []byte{0x02}
	RequestCountKey       = []byte{0x03}
	PendingRequestListKey = []byte{0x04}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func RequestStoreKey(requestID uint64) []byte {
	return append(RequestStoreKeyPrefix, sdk.Uint64ToBigEndian(requestID)...)
}

func ResultStoreKey(requestID uint64) []byte {
	return append(ResultStoreKeyPrefix, sdk.Uint64ToBigEndian(requestID)...)
}
