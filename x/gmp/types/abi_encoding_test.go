package types

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO: Make these tests more thorough
// Ref: https://github.com/ojo-network/ojo/issues/310
func TestEncodeABI(t *testing.T) {
	assetNameArray := [32]byte{}
	copy(assetNameArray[:], []byte("btc"))
	rates := [1]PriceFeedData{
		{
			AssetName:   assetNameArray,
			Value:       big.NewInt(50000),
			ResolveTime: big.NewInt(50000),
			Id:          big.NewInt(50000),
		},
	}
	disableResolve := false
	r, err := EncodeABI("postPrices", rates, disableResolve)
	require.NoError(t, err)
	fmt.Println(r)
}
