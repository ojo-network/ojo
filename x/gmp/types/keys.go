package types

const (
	// ModuleName is the name of the gmp module
	ModuleName = "gmp"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// RouterKey is the message route
	RouterKey = ModuleName

	// QuerierRoute is the query router key for the gmp module
	QuerierRoute = ModuleName
)

// KVStore key prefixes
var (
	ParamsKey = []byte{0x01}
)
