package tx

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	gasMultiplier = 3 / 2
)

// BroadcastTx attempts to generate, sign and broadcast a transaction with the
// given set of messages. It will also simulate gas requirements if necessary.
// It will return an error upon failure.
//
// Note, BroadcastTx is copied from the SDK except it removes a few unnecessary
// things like prompting for confirmation and printing the response. Instead,
// we return the TxResponse.
func BroadcastTx(clientCtx client.Context, txf tx.Factory, msgs ...sdk.Msg) (*sdk.TxResponse, error) {
	txf, err := txf.Prepare(clientCtx)
	if err != nil {
		return nil, err
	}

	_, adjusted, err := tx.CalculateGas(clientCtx, txf, msgs...)
	if err != nil {
		return nil, err
	}

	// make sure gas in enough to execute the txs
	txf = txf.WithGas(adjusted * gasMultiplier)

	unsignedTx, err := txf.BuildUnsignedTx(msgs...)
	if err != nil {
		return nil, err
	}

	unsignedTx.SetFeeGranter(clientCtx.GetFeeGranterAddress())

	if err = tx.Sign(clientCtx.CmdContext, txf, clientCtx.GetFromName(), unsignedTx, true); err != nil {
		return nil, err
	}

	txBytes, err := clientCtx.TxConfig.TxEncoder()(unsignedTx.GetTx())
	if err != nil {
		return nil, err
	}

	return clientCtx.BroadcastTx(txBytes)
}
