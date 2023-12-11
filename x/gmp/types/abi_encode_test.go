package types

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	message := GmpEncoder{
		PriceData: []PriceData{
			{
				AssetName:   [32]byte{},
				Price:       big.NewInt(1),
				ResolveTime: big.NewInt(1),
				MedianData: MedianData{
					BlockNums:  []*big.Int{big.NewInt(1)},
					Medians:    []*big.Int{big.NewInt(1)},
					Deviations: []*big.Int{big.NewInt(1)},
				},
			},
		},
		AssetNames:      [][32]byte{},
		ContractAddress: common.Address{},
		CommandSelector: [4]byte{},
		CommandParams:   []byte{},
	}

	bz, err := message.GMPEncode()
	require.NoError(t, err)
	vals, err := encoderSpec.Unpack(bz)
	require.NoError(t, err)

	require.Equal(t, message.AssetNames, vals[1].([][32]byte))
	require.Equal(t, message.ContractAddress, vals[2].(common.Address))
	require.Equal(t, message.CommandSelector, vals[3].([4]byte))
	require.Equal(t, message.CommandParams, vals[4].([]byte))
}
