package ante_test

import (
	"fmt"
	"testing"

	cmtrand "github.com/cometbft/cometbft/libs/rand"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/stretchr/testify/suite"

	ojoapp "github.com/ojo-network/ojo/app"
	appparams "github.com/ojo-network/ojo/app/params"
)

type IntegrationTestSuite struct {
	suite.Suite
	app       *ojoapp.App
	ctx       sdk.Context
	clientCtx client.Context
	txBuilder client.TxBuilder
}

func (s *IntegrationTestSuite) SetupTest() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(appparams.AccountAddressPrefix, appparams.AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(appparams.ValidatorAddressPrefix, appparams.ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(appparams.ConsNodeAddressPrefix, appparams.ConsNodePubKeyPrefix)

	app := ojoapp.Setup(s.T())
	ctx := app.BaseApp.NewContextLegacy(false, cmtproto.Header{
		ChainID: fmt.Sprintf("test-chain-%s", cmtrand.Str(4)),
		Height:  1,
	})

	s.app = app
	s.ctx = ctx

	encodingConfig := testutil.TestEncodingConfig{
		InterfaceRegistry: app.InterfaceRegistry(),
		Codec:             app.AppCodec(),
		TxConfig:          app.GetTxConfig(),
		Amino:             app.LegacyAmino(),
	}
	encodingConfig.Amino.RegisterConcrete(&testdata.TestMsg{}, "testdata.TestMsg", nil)
	testdata.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	s.clientCtx = client.Context{}.
		WithTxConfig(encodingConfig.TxConfig)
}

// CreateTestTx is a helper function to create a tx given multiple inputs.
func (suite *IntegrationTestSuite) CreateTestTx(privs []cryptotypes.PrivKey, accNums, accSeqs []uint64, chainID string) (xauthsigning.Tx, error) {
	var sigsV2 []signing.SignatureV2
	for i, priv := range privs {
		sigV2 := signing.SignatureV2{
			PubKey: priv.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode:  signing.SignMode(*suite.clientCtx.TxConfig.SignModeHandler().DefaultMode().Enum()),
				Signature: nil,
			},
			Sequence: accSeqs[i],
		}

		sigsV2 = append(sigsV2, sigV2)
	}
	err := suite.txBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return nil, err
	}

	sigsV2 = []signing.SignatureV2{}
	for i, priv := range privs {
		signerData := xauthsigning.SignerData{
			Address:       sdk.AccAddress(priv.PubKey().Bytes()).String(),
			ChainID:       chainID,
			AccountNumber: accNums[i],
			Sequence:      accSeqs[i],
			PubKey:        priv.PubKey(),
		}
		sigV2, err := tx.SignWithPrivKey(
			suite.ctx,
			signing.SignMode(*suite.clientCtx.TxConfig.SignModeHandler().DefaultMode().Enum()),
			signerData,
			suite.txBuilder,
			priv,
			suite.clientCtx.TxConfig,
			accSeqs[i],
		)
		if err != nil {
			return nil, err
		}

		sigsV2 = append(sigsV2, sigV2)
	}
	err = suite.txBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return nil, err
	}

	return suite.txBuilder.GetTx(), nil
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
