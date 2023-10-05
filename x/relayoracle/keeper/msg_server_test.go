package keeper_test

import (
	"github.com/ojo-network/ojo/x/relayoracle/types"
)

func (s *IntegrationTestSuite) TestMsgServer_UpdateGovParams() {
	govAccAddr := s.app.GovKeeper.GetGovernanceAccount(s.ctx).GetAddress().String()

	testCases := []struct {
		name      string
		req       *types.MsgGovUpdateParams
		expectErr bool
		errMsg    string
	}{
		{
			"valid max historical change",
			&types.MsgGovUpdateParams{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Keys:        []string{string(types.KeyMaxHistorical)},
				Changes: types.Params{
					MaxAllowedDenomsHistoricalQuery: 30,
				},
			},
			false,
			"",
		},
		{
			"valid max exchange query change",
			&types.MsgGovUpdateParams{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Keys:        []string{string(types.KeyMaxExchange)},
				Changes: types.Params{
					MaxAllowedDenomsExchangeQuery: 30,
				},
			},
			false,
			"",
		},
		{
			"valid packet timeout change",
			&types.MsgGovUpdateParams{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Keys:        []string{string(types.KeyPacketTimeout)},
				Changes: types.Params{
					PacketTimeout: 30,
				},
			},
			false,
			"",
		},
		{
			"valid ibc request enabled change",
			&types.MsgGovUpdateParams{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Keys:        []string{string(types.KeyIbcRequestEnabled)},
				Changes: types.Params{
					IbcRequestEnabled: true,
				},
			},
			false,
			"",
		},

		{
			"multiple valid ibc changes",
			&types.MsgGovUpdateParams{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Keys: []string{
					string(types.KeyIbcRequestEnabled),
					string(types.KeyPacketTimeout),
					string(types.KeyMaxExchange),
					string(types.KeyMaxHistorical),
				},
				Changes: types.Params{
					IbcRequestEnabled:               true,
					PacketTimeout:                   10,
					MaxAllowedDenomsHistoricalQuery: 11,
					MaxAllowedDenomsExchangeQuery:   12,
				},
			},
			false,
			"",
		},
		{
			"invalid key",
			&types.MsgGovUpdateParams{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Keys:        []string{"test"},
				Changes:     types.Params{},
			},
			true,
			"test is not a relay oracle param key",
		},

		{
			"bad authority",
			&types.MsgGovUpdateParams{
				Authority:   "ojo1zypqa76je7pxsdwkfah6mu9a583sju6xzthge3",
				Title:       "test",
				Description: "test",
				Keys:        []string{string(types.KeyIbcRequestEnabled)},
				Changes: types.Params{
					IbcRequestEnabled: true,
				},
			},
			true,
			"expected gov account as only signer for proposal message",
		},
	}

	for _, tc := range testCases {
		// reset params
		s.app.RelayOracle.SetParams(s.ctx, types.DefaultParams())

		s.Run(tc.name, func() {
			err := tc.req.ValidateBasic()
			if err == nil {
				_, err = s.msgServer.GovUpdateParams(s.ctx, tc.req)
			}

			if tc.expectErr {
				s.Require().ErrorContains(err, tc.errMsg)
			} else {
				s.Require().NoError(err)

				// check new params
				params := s.app.RelayOracle.GetParams(s.ctx)
				for _, key := range tc.req.Keys {
					s.T().Log(tc.name, params, tc.req.Changes)
					switch key {
					case string(types.KeyMaxHistorical):
						s.Require().EqualValues(params.MaxAllowedDenomsHistoricalQuery, tc.req.Changes.MaxAllowedDenomsHistoricalQuery)
					case string(types.KeyMaxExchange):
						s.Require().EqualValues(params.MaxAllowedDenomsExchangeQuery, tc.req.Changes.MaxAllowedDenomsExchangeQuery)
					case string(types.KeyPacketTimeout):
						s.Require().EqualValues(params.PacketTimeout, tc.req.Changes.PacketTimeout)
					case string(types.KeyIbcRequestEnabled):
						s.Require().EqualValues(params.IbcRequestEnabled, tc.req.Changes.IbcRequestEnabled)
					}
				}
			}
		})
	}
}
