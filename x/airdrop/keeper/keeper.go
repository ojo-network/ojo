package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey

	// the address capable of executing a MsgSetParams message. Typically, this
	// should be the x/gov module account.
	authority string
}

// NewKeeper constructs a new keeper for airdrop module.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	authority string,
) Keeper {
	return Keeper{
		cdc:       cdc,
		storeKey:  storeKey,
		authority: authority,
	}
}
