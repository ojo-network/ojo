package types

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToMap(t *testing.T) {
	rates := []ExchangeRate{
		{
			SymbolDenom: "btc",
			Rate:        big.NewInt(50000),
		},
	}
	r, err := EncodeExchangeRate(rates)
	require.NoError(t, err)
	fmt.Println(r)
}
