package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/ojo-network/ojo/x/relayoracle/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the relayoracle MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (ms msgServer) GovUpdateParams(
	goCtx context.Context,
	msg *types.MsgGovUpdateParams,
) (*types.MsgGovUpdateParamsResponse, error) {
	if ms.authority != msg.Authority {
		err := errors.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority; expected %s, got %s",
			ms.authority,
			msg.Authority,
		)
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	for _, key := range msg.Keys {
		switch key {
		case string(types.KeyIbcRequestEnabled):
			ms.SetIbcRequestEnabled(ctx, msg.Changes.IbcRequestEnabled)

		case string(types.KeyPacketTimeout):
			ms.SetPacketTimeout(ctx, msg.Changes.PacketTimeout)

		case string(types.KeyMaxExchange):
			ms.SetMaxQueryForExchangeRate(ctx, msg.Changes.MaxAllowedDenomsExchangeQuery)

		case string(types.KeyMaxHistorical):
			ms.SetMaxQueryForHistorical(ctx, msg.Changes.MaxAllowedDenomsHistoricalQuery)

		default:
			return nil, fmt.Errorf("%s is not a relay oracle param key", key)
		}
	}

	return &types.MsgGovUpdateParamsResponse{}, nil
}
