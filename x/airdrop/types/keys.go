package types

import (
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName is the name of the airdrop module
	ModuleName = "airdrop"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// RouterKey is the message route
	RouterKey = ModuleName

	// QuerierRoute is the query router key for the airdrop module
	QuerierRoute = ModuleName
)

// KVStore key prefixes
var (
	ParamsKey               = []byte{0x01}
	AirdropAccountKeyPrefix = []byte{0x02}
)

// AirdropAccountKey returns the store key for an airdrop account
func AirdropAccountKey(originAddress string, state AirdropAccount_State) (key []byte) {
	key = AirdropAccountKeyPrefix
	key = append(key, []byte(state.String())...)
	key = append(key, address.MustLengthPrefix([]byte(originAddress))...)
	return key
}

func AirdropIteratorKey(state AirdropAccount_State) (key []byte) {
	key = AirdropAccountKeyPrefix
	key = append(key, []byte(state.String())...)
	return key
}
