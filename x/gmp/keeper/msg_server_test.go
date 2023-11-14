package keeper_test

import (
	"github.com/cometbft/cometbft/crypto/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	appparams "github.com/ojo-network/ojo/app/params"
	"github.com/ojo-network/ojo/x/gmp/types"
)

var (
	pubKey    = secp256k1.GenPrivKey().PubKey()
	addr      = sdk.AccAddress(pubKey.Address())
	initCoins = sdk.NewCoins(sdk.NewCoin(appparams.BondDenom, sdk.NewInt(1000000000000000000)))
)

func (s *IntegrationTestSuite) TestMsgServer_SetParams() {
	gmpChannel := "channel-1"
	gmpAddress := "axelar1dv4u5k73pzqrxlzujxg3qp8kvc3pje7jtdvu72npnt5zhq05ejcsn5qme5"
	timeout := int64(1)
	SetParams(s, gmpAddress, gmpChannel, timeout)

	params := types.DefaultParams()

	s.Require().Equal(params, s.app.GmpKeeper.GetParams(s.ctx))
}

// SetParams sets the gmp module params
func SetParams(
	s *IntegrationTestSuite,
	gmpAddress string,
	gmpChannel string,
	gmpTimeout int64,
) {
	params := types.DefaultParams()
	params.GmpAddress = gmpAddress
	authority := s.app.GovKeeper.GetGovernanceAccount(s.ctx).GetAddress().String()

	msg := types.NewMsgSetParams(
		params.GmpAddress,
		params.GmpChannel,
		params.GmpTimeout,
		authority,
	)

	_, err := s.msgServer.SetParams(s.ctx, msg)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TestMsgServer_RelayPrices() {
	// Set default params
	app, ctx := s.app, s.ctx
	app.GmpKeeper.SetParams(ctx, types.DefaultParams())
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
			Amount: sdk.NewInt(1),
		},
		[]string{"BTC", "ATOM"},
		[]byte{1, 2, 3, 4},
		[]byte{1, 2, 3, 4},
		1000,
	)

	_, err := app.GmpKeeper.RelayPrice(ctx, msg)
	s.Require().NoError(err)
}
