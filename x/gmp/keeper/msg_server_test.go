package keeper_test

import (
	"github.com/ojo-network/ojo/x/gmp/types"
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
