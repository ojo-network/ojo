package tx

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gmptypes "github.com/ojo-network/ojo/x/gmp/types"
)

// TxCreatePayment creates a gmp payment transaction
func (c *Client) TxCreatePayment() (*sdk.TxResponse, error) {
	fromAddr, err := c.keyringRecord.GetAddress()
	if err != nil {
		return nil, err
	}

	msg := gmptypes.NewMsgCreatePayment(
		fromAddr.String(),
		"Arbitrum",
		"BTC",
		sdk.NewCoin("uojo", math.NewInt(10_000_000_000)),
		math.LegacyOneDec(),
		100,
	)

	return c.BroadcastTx(msg)
}
