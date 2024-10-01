package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	appparams "github.com/ojo-network/ojo/app/params"
	"github.com/ojo-network/ojo/x/gmp/types"
)

var (
	pubKey    = secp256k1.GenPrivKey().PubKey()
	addr      = sdk.AccAddress(pubKey.Address())
	initCoins = sdk.NewCoins(sdk.NewCoin(appparams.BondDenom, math.NewInt(1000000000000000000)))
)

func (s *IntegrationTestSuite) TestMsgServer_SetParams() {
	gmpChannel := "channel-1"
	gmpAddress := "axelar1dv4u5k73pzqrxlzujxg3qp8kvc3pje7jtdvu72npnt5zhq05ejcsn5qme5"
	timeout := int64(1000000)
	feeRecipient := "axelar1zl3rxpp70lmte2xr6c4lgske2fyuj3hupcsvcd"
	SetParams(s, gmpAddress, gmpChannel, timeout, feeRecipient)

	params := types.DefaultParams()

	s.Require().Equal(params, s.app.GmpKeeper.GetParams(s.ctx))
}

// SetParams sets the gmp module params
func SetParams(
	s *IntegrationTestSuite,
	gmpAddress string,
	gmpChannel string,
	gmpTimeout int64,
	feeRecipient string,
) {
	params := types.DefaultParams()
	params.GmpAddress = gmpAddress
	params.FeeRecipient = feeRecipient
	authority := s.app.GovKeeper.GetGovernanceAccount(s.ctx).GetAddress().String()

	msg := types.NewMsgSetParams(
		params.GmpAddress,
		params.GmpChannel,
		params.GmpTimeout,
		params.FeeRecipient,
		authority,
		params.DefaultGasEstimate,
	)

	_, err := s.msgServer.SetParams(s.ctx, msg)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TestMsgServer_RelayPrices() {
	// Set default params
	app, ctx := s.app, s.ctx
	s.Require().NoError(app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initCoins))
	s.Require().NoError(app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, initCoins))

	// Attempt a relay transaction
	msg := types.NewMsgRelay(
		addr.String(),
		"Ethereum",
		"0x0000",
		"0x0000",
		sdk.Coin{
			Denom:  "uojo",
			Amount: math.NewInt(1),
		},
		[]string{"BTC", "ATOM"},
		[]byte{1, 2, 3, 4},
		[]byte{1, 2, 3, 4},
		1000,
	)
	app.GmpKeeper.RelayPrice(ctx, msg)

	// Attempt a normal IBC transfer
	transferMsg := ibctransfertypes.NewMsgTransfer(
		ibctransfertypes.PortID,
		"channel-1",
		msg.Token,
		addr.String(),
		addr.String(),
		clienttypes.ZeroHeight(),
		uint64(1000),
		"memo",
	)
	app.TransferKeeper.Transfer(ctx, transferMsg)
}
