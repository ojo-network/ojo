package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/util/checkers"
)

var _ sdk.Msg = &MsgSetParams{}

func NewMsgSetParams(
	gmpAddress string,
	gmpChannel string,
	gmpTimeout int64,
	govAddress string,
) *MsgSetParams {
	params := &Params{
		GmpAddress: gmpAddress,
		GmpChannel: gmpChannel,
		GmpTimeout: gmpTimeout,
	}
	return &MsgSetParams{
		Params:    params,
		Authority: govAddress,
	}
}

// Type implements LegacyMsg interface
func (msg MsgSetParams) Type() string { return sdk.MsgTypeURL(&msg) }

// GetSignBytes implements sdk.Msg
func (msg MsgSetParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements sdk.Msg
func (msg MsgSetParams) GetSigners() []sdk.AccAddress {
	return checkers.Signers(msg.Authority)
}

// ValidateBasic Implements sdk.Msg
func (msg MsgSetParams) ValidateBasic() error {
	// TODO validate params
	return nil
}

func NewMsgRelay(
	relayer string,
	destinationChain string,
	destinationAddress string,
	token sdk.Coin,
	denoms []string,
) *MsgRelay {
	return &MsgRelay{
		Relayer:            relayer,
		DestinationChain:   destinationChain,
		DestinationAddress: destinationAddress,
		Token:              token,
		Denoms:             denoms,
	}
}

// Type implements LegacyMsg interface
func (msg MsgRelay) Type() string { return sdk.MsgTypeURL(&msg) }

// GetSignBytes implements sdk.Msg
func (msg MsgRelay) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements sdk.Msg
func (msg MsgRelay) GetSigners() []sdk.AccAddress {
	return checkers.Signers(msg.Relayer)
}

// ValidateBasic Implements sdk.Msg
func (msg MsgRelay) ValidateBasic() error {
	// TODO validate relay msg
	return nil
}
