package tx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

// TxDelegateFeedConsent delegates the feed consent to the given feeder
// from the OjoClient's validator
func (c *Client) TxDelegateFeedConsent(
	feeder sdk.AccAddress,
) (*sdk.TxResponse, error) {
	addr, err := c.keyringRecord.GetAddress()
	if err != nil {
		return nil, err
	}

	validator := sdk.ValAddress(addr)

	msg := oracletypes.NewMsgDelegateFeedConsent(
		validator,
		feeder,
	)
	return c.BroadcastTx(msg)
}
