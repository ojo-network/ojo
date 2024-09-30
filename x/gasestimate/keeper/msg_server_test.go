package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	appparams "github.com/ojo-network/ojo/app/params"
	"github.com/ojo-network/ojo/x/gasestimate/types"
)

var (
	pubKey    = secp256k1.GenPrivKey().PubKey()
	addr      = sdk.AccAddress(pubKey.Address())
	initCoins = sdk.NewCoins(sdk.NewCoin(appparams.BondDenom, math.NewInt(1000000000000000000)))
)

func (s *IntegrationTestSuite) TestMsgServer_SetParams() {
	gasestimateChannel := "channel-1"
	gasestimateAddress := "axelar1dv4u5k73pzqrxlzujxg3qp8kvc3pje7jtdvu72npnt5zhq05ejcsn5qme5"
	timeout := int64(1)
	feeRecipient := "axelar1zl3rxpp70lmte2xr6c4lgske2fyuj3hupcsvcd"
	SetParams(s, gasestimateAddress, gasestimateChannel, timeout, feeRecipient)

	params := types.DefaultParams()

	s.Require().Equal(params, s.app.GasEstimateKeeper.GetParams(s.ctx))
}

// SetParams sets the gasestimate module params
func SetParams(
	s *IntegrationTestSuite,
	gasestimateAddress string,
	gasestimateChannel string,
	gasestimateTimeout int64,
	feeRecipient string,
) {
	params := types.DefaultParams()
	authority := s.app.GovKeeper.GetGovernanceAccount(s.ctx).GetAddress().String()

	msg := types.NewMsgSetParams(
		params.ContractRegistry,
		authority,
	)

	_, err := s.msgServer.SetParams(s.ctx, msg)
	s.Require().NoError(err)
}
