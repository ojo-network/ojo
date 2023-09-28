package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"gopkg.in/yaml.v3"

	"github.com/ojo-network/ojo/util/checkers"
)

var (
	_ sdk.Msg = &MsgGovUpdateParams{}
)

// NewMsgUpdateParams will creates a new MsgUpdateParams instance
func NewMsgUpdateParams(authority, title, description string, keys []string, changes Params) *MsgGovUpdateParams {
	return &MsgGovUpdateParams{
		Title:       title,
		Description: description,
		Authority:   authority,
		Keys:        keys,
		Changes:     changes,
	}
}

// Type implements Msg interface
func (msg MsgGovUpdateParams) Type() string { return sdk.MsgTypeURL(&msg) }

// String implements the Stringer interface.
func (msg MsgGovUpdateParams) String() string {
	out, _ := yaml.Marshal(msg)
	return string(out)
}

// GetSignBytes implements Msg
func (msg MsgGovUpdateParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgGovUpdateParams) GetSigners() []sdk.AccAddress {
	return checkers.Signers(msg.Authority)
}

// ValidateBasic implements Msg and validates params for each param key
// specified in the proposal. If one param is invalid, the whole proposal
// will fail to go through.
func (msg MsgGovUpdateParams) ValidateBasic() error {
	if err := checkers.ValidateProposal(msg.Title, msg.Description, msg.Authority); err != nil {
		return err
	}

	for _, key := range msg.Keys {
		switch key {
		case string(KeyIbcRequestEnabled):
			return validateBool(msg.Changes.IbcRequestEnabled)

		case string(KeyPacketTimeout):
			return validateUint64(msg.Changes.PacketTimeout)

		default:
			return fmt.Errorf("%s is not a relay oracle param key", key)
		}
	}

	return nil
}
