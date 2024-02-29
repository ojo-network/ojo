package tx

import (
	"os"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	cmtjsonclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ojoapp "github.com/ojo-network/ojo/app"
)

const (
	gasAdjustment = 1
)

// TxClient is a wrapper around the cosmos sdk client context and transaction
// factory for signing and broadcasting transactions
type Client struct {
	ChainID        string
	CMTRPCEndpoint string

	ClientContext *client.Context

	keyringKeyring keyring.Keyring
	keyringRecord  *keyring.Record
	txFactory      *tx.Factory
}

// Initializes a cosmos sdk client context and transaction factory for
// signing and broadcasting transactions
func NewClient(
	chainID string,
	cmtrpcEndpoint string,
	accountName string,
	accountMnemonic string,
) (c *Client, err error) {
	c = &Client{
		ChainID:        chainID,
		CMTRPCEndpoint: cmtrpcEndpoint,
	}

	c.keyringRecord, c.keyringKeyring, err = CreateAccountFromMnemonic(accountName, accountMnemonic)
	if err != nil {
		return nil, err
	}

	err = c.createClientContext()
	if err != nil {
		return nil, err
	}
	c.createTxFactory()

	return c, err
}

func (c *Client) createClientContext() error {
	encoding := ojoapp.MakeEncodingConfig()
	fromAddress, _ := c.keyringRecord.GetAddress()

	cmtHTTPClient, err := cmtjsonclient.DefaultHTTPClient(c.CMTRPCEndpoint)
	if err != nil {
		return err
	}

	cmtRPCClient, err := rpchttp.NewWithClient(c.CMTRPCEndpoint, "/websocket", cmtHTTPClient)
	if err != nil {
		return err
	}

	c.ClientContext = &client.Context{
		ChainID:           c.ChainID,
		InterfaceRegistry: encoding.InterfaceRegistry,
		Output:            os.Stderr,
		BroadcastMode:     flags.BroadcastSync,
		TxConfig:          encoding.TxConfig,
		AccountRetriever:  authtypes.AccountRetriever{},
		Codec:             encoding.Codec,
		LegacyAmino:       encoding.Amino,
		Input:             os.Stdin,
		NodeURI:           c.CMTRPCEndpoint,
		Client:            cmtRPCClient,
		Keyring:           c.keyringKeyring,
		FromAddress:       fromAddress,
		FromName:          c.keyringRecord.Name,
		From:              c.keyringRecord.Name,
		OutputFormat:      "json",
		UseLedger:         false,
		Simulate:          false,
		GenerateOnly:      false,
		Offline:           false,
		SkipConfirm:       true,
	}

	return nil
}

func (c *Client) createTxFactory() {
	factory := tx.Factory{}.
		WithAccountRetriever(c.ClientContext.AccountRetriever).
		WithChainID(c.ChainID).
		WithTxConfig(c.ClientContext.TxConfig).
		WithGasAdjustment(gasAdjustment).
		WithKeybase(c.ClientContext.Keyring).
		WithSignMode(signing.SignMode_SIGN_MODE_DIRECT).
		WithSimulateAndExecute(true)
	c.txFactory = &factory
}

func (c *Client) Address() (sdk.AccAddress, error) {
	return c.keyringRecord.GetAddress()
}

func (c *Client) BroadcastTx(msgs ...sdk.Msg) (*sdk.TxResponse, error) {
	return BroadcastTx(*c.ClientContext, *c.txFactory, msgs...)
}
