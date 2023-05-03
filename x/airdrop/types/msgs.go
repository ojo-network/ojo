package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/util/checkers"
)

var (
	_ sdk.Msg = &MsgSetParams{}
	_ sdk.Msg = &MsgCreateAirdropAccount{}
	_ sdk.Msg = &MsgClaimAirdrop{}
)

func NewMsgSetParams(
	expiryBlock uint64,
	delegationRequirement *sdk.Dec,
	airdropFactor *sdk.Dec,
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

func NewMsgCreateAirdropAccount(
	address string,
	tokensToReceive *sdk.DecCoin,
	vestingEndTime int64,
) *MsgCreateAirdropAccount {
	return &MsgCreateAirdropAccount{
		Address:         address,
		TokensToReceive: tokensToReceive,
		VestingEndTime:  vestingEndTime,
	}
}

// Type implements LegacyMsg interface
func (msg MsgCreateAirdropAccount) Type() string { return sdk.MsgTypeURL(&msg) }

// GetSignBytes implements sdk.Msg
func (msg MsgCreateAirdropAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements sdk.Msg
func (msg MsgCreateAirdropAccount) GetSigners() []sdk.AccAddress {
	return checkers.Signers()
}

// ValidateBasic Implements sdk.Msg
func (msg MsgCreateAirdropAccount) ValidateBasic() error {
	// TODO validate address
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

// GetSignBytes implements sdk.Msg
func (msg MsgClaimAirdrop) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements sdk.Msg
func (msg MsgClaimAirdrop) GetSigners() []sdk.AccAddress {
	return checkers.Signers(msg.FromAddress)
}

// ValidateBasic implements sdk.Msg
func (msg MsgClaimAirdrop) ValidateBasic() error {
	// TODO validate addresses
	return nil
}
