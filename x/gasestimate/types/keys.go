package types

const (
	// ModuleName is the name of the gasestimate module
	ModuleName = "gasestimate"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// RouterKey is the message route
	RouterKey = ModuleName

	// QuerierRoute is the query router key for the gasestimate module
	QuerierRoute = ModuleName
)

// KVStore key prefixes
var (
	ParamsKey      = []byte{0x01}
	GasEstimateKey = []byte{0x02}
)

// KeyPrefixGasEstimates is the prefix for the gas estimates
func KeyPrefixGasEstimate(network string) []byte {
	return append(GasEstimateKey, []byte(network)...)
}
