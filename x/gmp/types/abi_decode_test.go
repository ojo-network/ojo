package types

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// TestGmpData tests the GmpData struct by encoding and decoding it.
func TestGmpData(t *testing.T) {
	g := GmpDecoder{
		AssetNames:      [][32]byte{{1}},
		ContractAddress: common.HexToAddress("0x0000001"),
		CommandSelector: [4]byte{1},
		CommandParams:   []byte{1},
		Timestamp:       big.NewInt(1),
	}
	payload, err := decoderSpec.Pack(
		g.AssetNames,
		g.ContractAddress,
		g.CommandSelector,
		g.CommandParams,
		g.Timestamp,
	)
	require.NoError(t, err)
	newGmpData, err := NewGmpDecoder(payload)
	require.NoError(t, err)

	require.Equal(t, g, newGmpData)
}

func TestGetDenoms(t *testing.T) {
	assetNamesString := []string{
		"BTC",
		"ETH",
		"USDT",
		"BNB",
		"ADA",
		"FOOBARFOOBARFOOBARFOOBARFOOBARFO", // maxmimum allowed length
	}
	assetNamesAsBytes := make([][32]byte, len(assetNamesString))
	for i, name := range assetNamesString {
		copy(assetNamesAsBytes[i][:], name)
	}

	g := GmpDecoder{
		AssetNames: assetNamesAsBytes,
	}
	denoms := g.GetDenoms()
	require.Equal(t, assetNamesString, denoms)
}
