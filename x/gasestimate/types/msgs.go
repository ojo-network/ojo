package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/util/checkers"
)

var _ sdk.Msg = &MsgSetParams{}

func NewMsgSetParams(
	contractRegistry []*Contract,
	govAddress string,
) *MsgSetParams {
	params := &Params{
		ContractRegistry: contractRegistry,
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
