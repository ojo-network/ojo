package types

import "github.com/cosmos/cosmos-sdk/types/address"

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
	ParamsKey        = []byte{0x01}
	PaymentKeyPrefix = []byte{0x02}
)

// PaymentKey returns the store key for a payment
func PaymentKey(originAddress string, denom string) (key []byte) {
	key = PaymentKeyPrefix
	key = append(key, []byte(denom)...)
	key = append(key, address.MustLengthPrefix([]byte(originAddress))...)
	return key
}
