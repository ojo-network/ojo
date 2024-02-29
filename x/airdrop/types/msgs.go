package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/util/checkers"
)

const (
	TypeMsgClaimAirdrop = "claim_airdrop"
)

var (
	_ sdk.Msg = &MsgSetParams{}
	_ sdk.Msg = &MsgClaimAirdrop{}
)

func NewMsgSetParams(
	expiryBlock uint64,
	delegationRequirement *math.LegacyDec,
	airdropFactor *math.LegacyDec,
	govAddress string,
) *MsgSetParams {
	params := &Params{
		ExpiryBlock:           expiryBlock,
		DelegationRequirement: delegationRequirement,
		AirdropFactor:         airdropFactor,
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

func NewMsgClaimAirdrop(
	fromAddress string,
	toAddress string,
) *MsgClaimAirdrop {
	return &MsgClaimAirdrop{
		FromAddress: fromAddress,
		ToAddress:   toAddress,
	}
}

// Type implements LegacyMsg interface
func (msg MsgClaimAirdrop) Type() string { return sdk.MsgTypeURL(&msg) }

// GetSigners implements sdk.Msg
func (msg MsgClaimAirdrop) GetSigners() []sdk.AccAddress {
	return checkers.Signers(msg.FromAddress)
}

// ValidateBasic implements sdk.Msg
func (msg MsgClaimAirdrop) ValidateBasic() error {
	// TODO validate addresses
	return nil
}
