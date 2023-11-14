package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgRelayPriceValidateBasic(t *testing.T) {
	coins, err := sdk.ParseCoinNormalized("100uojo")
	require.NoError(t, err)
	price := NewMsgRelay(
		"relayer",
		"axelar-1",
		"axelar1dv4u5k73pzqrxlzujxg3qp8kvc3pje7jtdvu72npnt5zhq05ejcsn5qme5",
		coins,
		[]string{"uojo"},
		[]byte("1234"),
		[]byte("command-params"),
		100,
	)
	err = price.ValidateBasic()
	require.NoError(t, err)
	price.CommandSelector = []byte("12345")
	err = price.ValidateBasic()
	require.Error(t, err)
}
