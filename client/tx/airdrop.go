package tx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	airdroptypes "github.com/ojo-network/ojo/x/airdrop/types"
)

func (c *Client) TxClaimAirdrop(
	fromAddress string,
	toAddress string,
) (*sdk.TxResponse, error) {
	msg := airdroptypes.NewMsgClaimAirdrop(
		fromAddress,
		toAddress,
	)
	return c.BroadcastTx(msg)
}
