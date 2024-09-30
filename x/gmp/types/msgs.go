package types

import (
	fmt "fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/util/checkers"
)

var _ sdk.Msg = &MsgSetParams{}

func NewMsgSetParams(
	gmpAddress string,
	gmpChannel string,
	gmpTimeout int64,
	feeRecipient string,
	govAddress string,
) *MsgSetParams {
	params := &Params{
		GmpAddress:   gmpAddress,
		GmpChannel:   gmpChannel,
		GmpTimeout:   gmpTimeout,
		FeeRecipient: feeRecipient,
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

// GetSigners implements sdk.Msg
func (msg MsgRelayPrice) GetSigners() []sdk.AccAddress {
	return checkers.Signers(msg.Relayer)
}

// ValidateBasic Implements sdk.Msg
func (msg MsgRelayPrice) ValidateBasic() error {
	if len(msg.Denoms) == 0 {
		return fmt.Errorf("denoms cannot be empty")
	}
	if msg.Timestamp == 0 {
		return fmt.Errorf("timestamp cannot be empty")
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

func NewMsgCreatePayment(
	relayer string,
	destinationChain string,
	denom string,
	token sdk.Coin,
	deviation math.LegacyDec,
	heartbeat int64,
) *MsgCreatePayment {
	return &MsgCreatePayment{
		Relayer: relayer,
		Payment: &Payment{
			Relayer:          relayer,
			DestinationChain: destinationChain,
			Denom:            denom,
			Token:            token,
			Deviation:        deviation,
			Heartbeat:        heartbeat,
		},
	}
}

// Type implements LegacyMsg interface
func (msg MsgCreatePayment) Type() string { return sdk.MsgTypeURL(&msg) }

// GetSigners implements sdk.Msg
func (msg MsgRelayPrice) MsgCreatePayment() []sdk.AccAddress {
	return checkers.Signers(msg.Relayer)
}

// ValidateBasic Implements sdk.Msg
func (msg MsgCreatePayment) ValidateBasic() error {
	if len(msg.Relayer) == 0 {
		return fmt.Errorf("relayer cannot be empty")
	}
	if msg.Payment.DestinationChain == "" {
		return fmt.Errorf("destinationChain cannot be empty")
	}
	if msg.Payment.Denom == "" {
		return fmt.Errorf("denom cannot be empty")
	}
	if msg.Payment.Token.IsZero() {
		return fmt.Errorf("token cannot be zero")
	}
	// deviation must be between 0.5 and 50
	if msg.Payment.Deviation.LT(math.LegacyNewDecWithPrec(5, 1)) || msg.Payment.Deviation.GT(math.LegacyNewDec(50)) {
		return fmt.Errorf("deviation must be between 0.5 and 50")
	}
	if msg.Payment.Heartbeat <= 0 {
		return fmt.Errorf("heartbeat must be greater than 0")
	}

	return nil
}
