package gmpmiddleware

import (
	"math/big"
	"testing"

	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	"github.com/ethereum/go-ethereum/common"
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

// TestGmpData tests the GmpData struct by encoding and decoding it.
func TestGmpData(t *testing.T) {
	gmpData := GmpData{
		AssetNames:      [][32]byte{{1}},
		ContractAddress: common.HexToAddress("0x0000001"),
		CommandSelector: [4]byte{1},
		CommandParams:   []byte{1},
		Timestamp:       big.NewInt(1),
	}
	payload, err := gmpData.Encode()
	require.NoError(t, err)
	newGmpData, err := NewGmpData(payload)
	require.NoError(t, err)

	require.Equal(t, gmpData, newGmpData)
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
