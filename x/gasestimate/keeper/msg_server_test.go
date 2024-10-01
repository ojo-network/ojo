package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	appparams "github.com/ojo-network/ojo/app/params"
	"github.com/ojo-network/ojo/x/gasestimate/types"
)

var (
	pubKey           = secp256k1.GenPrivKey().PubKey()
	addr             = sdk.AccAddress(pubKey.Address())
	initCoins        = sdk.NewCoins(sdk.NewCoin(appparams.BondDenom, math.NewInt(1000000000000000000)))
	contractRegistry = []*types.Contract{
		{
			Address: "0x5BB3E85f91D08fe92a3D123EE35050b763D6E6A7",
			Network: "Ethereum",
		},
	}
)

func (s *IntegrationTestSuite) TestMsgServer_SetParams() {
	SetParams(s, contractRegistry, "1000000", "1.5")

	params := types.DefaultParams()

	s.Require().Equal(params, s.app.GasEstimateKeeper.GetParams(s.ctx))
}

// SetParams sets the gasestimate module params
func SetParams(
	s *IntegrationTestSuite,
	contractRegistry []*types.Contract,
	gasLimit string,
	gasAdjustment string,
) {
	params := types.DefaultParams()
	authority := s.app.GovKeeper.GetGovernanceAccount(s.ctx).GetAddress().String()

	msg := types.NewMsgSetParams(
		params.ContractRegistry,
		authority,
		gasLimit,
		gasAdjustment,
	)

	_, err := s.msgServer.SetParams(s.ctx, msg)
	s.Require().NoError(err)
}
