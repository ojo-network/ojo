package types

import (
	"github.com/ojo-network/ojo/util"
)

const (
	// ModuleName is the name of the symbiotic module
	ModuleName = "symbiotic"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// RouterKey is the message route
	RouterKey = ModuleName

	// QuerierRoute is the query router key for the symbiotic module
	QuerierRoute = ModuleName
)

// KVStore key prefixes
var (
	ParamsKey              = []byte{0x01}
	CachedBlockHashPrefix  = []byte{0x02}
	CachedHeaderInfoPrefix = []byte{0x03}
)

// CachedBlockHashKey returns the store key for a CachedBlockHash
func CachedBlockHashKey(blockHeight uint64) (key []byte) {
	return util.ConcatBytes(0, CachedBlockHashPrefix, util.UintWithNullPrefix(blockHeight))
}

// CachedHeaderInfoKey returns the store key for a CachedHeaderInfo
func CachedHeaderInfoKey(blockHeight uint64) (key []byte) {
	return util.ConcatBytes(0, CachedHeaderInfoPrefix, util.UintWithNullPrefix(blockHeight))
}
