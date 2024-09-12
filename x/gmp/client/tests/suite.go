package tests

import (
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	"github.com/ojo-network/ojo/x/gmp/client/cli"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func NewIntegrationTestSuite(cfg network.Config) *IntegrationTestSuite {
	return &IntegrationTestSuite{cfg: cfg}
}

func (s *IntegrationTestSuite) SetupSuite() {
	t := s.T()
	t.Log("setting up integration test suite")

	var err error
	s.network, err = network.New(t, t.TempDir(), s.cfg)
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")

	s.network.Cleanup()
}

// TODO: Fix tx raw log not having "no gmp account found" message in it.
// Ref: https://github.com/ojo-network/ojo/issues/308
func (s *IntegrationTestSuite) TestRelayGmp() {
	s.T().Skip()

	val := s.network.Validators[0]

	testCases := []struct {
		name         string
		args         []string
		expectErr    bool
		expectedCode uint32
		respType     proto.Message
	}{
		{
			name: "invalid from address",
			args: []string{
				"foo",
				s.network.Validators[1].Address.String(),
			},
			expectErr: true,
			respType:  &sdk.TxResponse{},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cli.GetCmdRelay(), tc.args)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())

				txResp := tc.respType.(*sdk.TxResponse)
				s.Require().Contains(txResp.RawLog, "unable to relay to gmp")
			}
		})
	}
}
