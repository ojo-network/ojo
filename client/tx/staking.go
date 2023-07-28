package tx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (c *Client) TxDelegate(
	fromAddress sdk.AccAddress,
	validatorAddress sdk.ValAddress,
	amount sdk.Coin,
) (*sdk.TxResponse, error) {
	msg := stakingtypes.NewMsgDelegate(
		fromAddress,
		validatorAddress,
		amount,
	)
	return c.BroadcastTx(msg)
}
