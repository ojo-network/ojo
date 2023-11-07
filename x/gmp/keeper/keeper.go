package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/ojo-network/ojo/app/ibctransfer"
	"github.com/ojo-network/ojo/x/gmp/types"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	storeKey     storetypes.StoreKey
	oracleKeeper types.OracleKeeper
	ibcKeeper    ibctransfer.Keeper
	// the address capable of executing a MsgSetParams message. Typically, this
	// should be the x/gov module account.
	authority string
}

// NewKeeper constructs a new keeper for gmp module.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	oracleKeeper types.OracleKeeper,
	ibcKeeper ibctransfer.Keeper,
	authority string,
) Keeper {
	return Keeper{
		cdc:          cdc,
		storeKey:     storeKey,
		authority:    authority,
		oracleKeeper: oracleKeeper,
		ibcKeeper:    ibcKeeper,
	}
}
