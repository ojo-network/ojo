package types

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
	ParamsKey        = []byte{0x01}
	PaymentKeyPrefix = []byte{0x02}
)
