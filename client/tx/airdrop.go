package tx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	airdroptypes "github.com/ojo-network/ojo/x/airdrop/types"
)

func (c *Client) TxSubmitAirdropProposal(
	params *airdroptypes.Params,
) (*sdk.TxResponse, error) {

	msg := airdroptypes.NewMsgSetParams(
		params.ExpiryBlock,
		params.DelegationRequirement,
		params.AirdropFactor,
	)

	address, err := c.keyringRecord.GetAddress()
	if err != nil {
		return nil, err
	}

	proposalMessage, err := govtypesv1.NewMsgSubmitProposal(
		[]sdk.Msg{msg},
		sdk.NewCoins(sdk.NewCoin("uojo", sdk.NewInt(10000000))),
		address.String(),
		"",
	)
	if err != nil {
		return nil, err
	}

	return c.BroadcastTx(proposalMessage)
}
