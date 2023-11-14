package gmpmiddleware

import (
	"testing"

	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	"github.com/ojo-network/ojo/x/gmp/types"
	"github.com/stretchr/testify/require"
)

func TestVerifyParams(t *testing.T) {
	params := types.Params{
		GmpAddress: "gmpAddress",
		GmpChannel: "gmpChannel",
	}
	err := verifyParams(params, "gmpAddress", "gmpChannel")
	require.NoError(t, err)

	err = verifyParams(params, "notAddress", "notChannel")
	require.Error(t, err)
}

func TestParseDenom(t *testing.T) {
	packet := channeltypes.Packet{
		SourcePort:    "ibc",
		SourceChannel: "channel-0",
	}
	denom := "ibc/1D2F3A4B5C6D7E8F9A0B1C2D3E4F5A6B7C8D9E0F"
	parsedDenom := parseDenom(packet, denom)
	require.Equal(t, "//ibc/1D2F3A4B5C6D7E8F9A0B1C2D3E4F5A6B7C8D9E0F", parsedDenom)

	denom = "ibc/1D2F3A4B5C6D7E8F9A0B1C2D3E4F5A6B7C8D9E0F/ibc/1D2F3A4B5C6D7E8F9A0B1C2D3E4F5A6B7C8D9E0F"
	parsedDenom = parseDenom(packet, denom)
	require.Equal(t, "//ibc/1D2F3A4B5C6D7E8F9A0B1C2D3E4F5A6B7C8D9E0F/ibc/1D2F3A4B5C6D7E8F9A0B1C2D3E4F5A6B7C8D9E0F", parsedDenom)
}
