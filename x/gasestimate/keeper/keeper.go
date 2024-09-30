package keeper

import (
	"fmt"

	"cosmossdk.io/log"

	sdk "github.com/cosmos/cosmos-sdk/types"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ojo-network/ojo/app/ibctransfer"
	"github.com/ojo-network/ojo/x/gasestimate/types"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	IBCKeeper  *ibctransfer.Keeper
	BankKeeper types.BankKeeper
	// the address capable of executing a MsgSetParams message. Typically, this
	// should be the x/gov module account.
	authority string
}

// NewKeeper constructs a new keeper for gasestimate module.
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

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// set gas estimates
func (k Keeper) SetGasEstimate(ctx sdk.Context, gasEstimate types.GasEstimate) {
	store := ctx.KVStore(k.storeKey)
	gasEstimateBz, err := k.cdc.Marshal(&gasEstimate)
	if err != nil {
		panic(err)
	}
	store.Set(types.KeyPrefixGasEstimate(gasEstimate.Network), gasEstimateBz)
}

// get gas estimates
func (k Keeper) GetGasEstimate(ctx sdk.Context, network string) (types.GasEstimate, error) {
	store := ctx.KVStore(k.storeKey)
	gasEstimateBz := store.Get(types.KeyPrefixGasEstimate(network))
	if gasEstimateBz == nil {
		return types.GasEstimate{}, fmt.Errorf("gas estimate not found for network: %s", network)
	}
	var gasEstimate types.GasEstimate
	err := k.cdc.Unmarshal(gasEstimateBz, &gasEstimate)
	if err != nil {
		return types.GasEstimate{}, err
	}
	return gasEstimate, nil
}
