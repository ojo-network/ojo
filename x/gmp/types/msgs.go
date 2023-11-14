package types

import (
	fmt "fmt"

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
	ojoContractAddress string,
	clientContractAddress string,
	token sdk.Coin,
	denoms []string,
	commandSelector []byte,
	commandParams []byte,
	timestamp int64,
) *MsgRelayPrice {
	return &MsgRelayPrice{
		Relayer:               relayer,
		DestinationChain:      destinationChain,
		ClientContractAddress: clientContractAddress,
		OjoContractAddress:    ojoContractAddress,
		Token:                 token,
		Denoms:                denoms,
		CommandSelector:       commandSelector,
		CommandParams:         commandParams,
		Timestamp:             timestamp,
	}
}

// Type implements LegacyMsg interface
func (msg MsgRelayPrice) Type() string { return sdk.MsgTypeURL(&msg) }

// GetSignBytes implements sdk.Msg
func (msg MsgRelayPrice) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements sdk.Msg
func (msg MsgRelayPrice) GetSigners() []sdk.AccAddress {
	return checkers.Signers(msg.Relayer)
}

// ValidateBasic Implements sdk.Msg
func (msg MsgRelayPrice) ValidateBasic() error {
	if len(msg.CommandParams) == 0 {
		return fmt.Errorf("commandParams cannot be empty")
	}
	if len(msg.CommandSelector) != 4 {
		return fmt.Errorf("commandSelector length must be 4")
	}
	if len(msg.Denoms) == 0 {
		return fmt.Errorf("denoms cannot be empty")
	}
	if msg.Timestamp == 0 {
		return fmt.Errorf("timestamp cannot be empty")
	}
	if msg.ClientContractAddress == "" {
		return fmt.Errorf("clientContractAddress cannot be empty")
	}
	if msg.OjoContractAddress == "" {
		return fmt.Errorf("cjoContractAddress cannot be empty")
	}
	if msg.DestinationChain == "" {
		return fmt.Errorf("destinationChain cannot be empty")
	}
	if msg.Relayer == "" {
		return fmt.Errorf("relayer cannot be empty")
	}

	// Make sure no denoms are above 32 bytes
	for _, denom := range msg.Denoms {
		if len(denom) > 32 {
			return fmt.Errorf("denom %s is too long", denom)
		}
	}

	return nil
}
