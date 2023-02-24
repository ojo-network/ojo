package tx

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ojoapp "github.com/ojo-network/ojo/app"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	tmjsonclient "github.com/tendermint/tendermint/rpc/jsonrpc/client"
)

const (
	gasAdjustment = 1
)

// TxClient is a wrapper around the cosmos sdk client context and transaction
// factory for signing and broadcasting transactions
type TxClient struct {
	ChainID       string
	TMRPCEndpoint string

	ClientContext *client.Context

	keyringKeyring keyring.Keyring
	keyringRecord  *keyring.Record
	txFactory      *tx.Factory
}

// Initializes a cosmos sdk client context and transaction factory for
// signing and broadcasting transactions
func NewTxClient(
	chainID string,
	tmrpcEndpoint string,
	accountName string,
	accountMnemonic string,
) (tc *TxClient, err error) {
	tc = &TxClient{
		ChainID:       chainID,
		TMRPCEndpoint: tmrpcEndpoint,
	}

	tc.keyringRecord, tc.keyringKeyring, err = CreateAccountFromMnemonic(accountName, accountMnemonic)
	if err != nil {
		return nil, err
	}

	err = tc.createClientContext()
	if err != nil {
		return nil, err
	}
	tc.createTxFactory()

	return tc, err
}

func (tc *TxClient) createClientContext() error {
	encoding := ojoapp.MakeEncodingConfig()
	fromAddress, _ := tc.keyringRecord.GetAddress()

	tmHTTPClient, err := tmjsonclient.DefaultHTTPClient(tc.TMRPCEndpoint)
	if err != nil {
		return err
	}

	tmRPCClient, err := rpchttp.NewWithClient(tc.TMRPCEndpoint, "/websocket", tmHTTPClient)
	if err != nil {
		return err
	}

	tc.ClientContext = &client.Context{
		ChainID:           tc.ChainID,
		InterfaceRegistry: encoding.InterfaceRegistry,
		Output:            os.Stderr,
		BroadcastMode:     flags.BroadcastBlock,
		TxConfig:          encoding.TxConfig,
		AccountRetriever:  authtypes.AccountRetriever{},
		Codec:             encoding.Codec,
		LegacyAmino:       encoding.Amino,
		Input:             os.Stdin,
		NodeURI:           tc.TMRPCEndpoint,
		Client:            tmRPCClient,
		Keyring:           tc.keyringKeyring,
		FromAddress:       fromAddress,
		FromName:          tc.keyringRecord.Name,
		From:              tc.keyringRecord.Name,
		OutputFormat:      "json",
		UseLedger:         false,
		Simulate:          false,
		GenerateOnly:      false,
		Offline:           false,
		SkipConfirm:       true,
	}
	return nil
}

func (tc *TxClient) createTxFactory() {
	factory := tx.Factory{}.
		WithAccountRetriever(tc.ClientContext.AccountRetriever).
		WithChainID(tc.ChainID).
		WithTxConfig(tc.ClientContext.TxConfig).
		WithGasAdjustment(gasAdjustment).
		WithKeybase(tc.ClientContext.Keyring).
		WithSignMode(signing.SignMode_SIGN_MODE_DIRECT).
		WithSimulateAndExecute(true)
	tc.txFactory = &factory
}

func (c *TxClient) BroadcastTx(msgs ...sdk.Msg) (*sdk.TxResponse, error) {
	return BroadcastTx(*c.ClientContext, *c.txFactory, msgs...)
}
