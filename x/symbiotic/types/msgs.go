package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/util/checkers"
)

var _ sdk.Msg = &MsgSetParams{}

func NewMsgSetParams(
	govAddress string,
	middlewareAddress string,
	symbioticSyncPeriod int64,
	maximumCachedBlockHashes uint64,
) *MsgSetParams {
	params := &Params{
		MiddlewareAddress:        middlewareAddress,
		SymbioticSyncPeriod:      symbioticSyncPeriod,
		MaximumCachedBlockHashes: maximumCachedBlockHashes,
	}
	return &MsgSetParams{
		Params:    params,
		Authority: govAddress,
	}
}

// Type implements LegacyMsg interface
func (msg MsgSetParams) Type() string { return sdk.MsgTypeURL(&msg) }

// GetSigners implements sdk.Msg
func (msg MsgSetParams) GetSigners() []sdk.AccAddress {
	return checkers.Signers(msg.Authority)
}

// ValidateBasic Implements sdk.Msg
func (msg MsgSetParams) ValidateBasic() error {
	// TODO validate params
	return nil
}
