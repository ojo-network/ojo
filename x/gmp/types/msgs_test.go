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
		"0x010",
		"0x020",
		coins,
		[]string{"uojo"},
		[]byte("1234"),
		[]byte("command-params"),
		100,
	)
	err = price.ValidateBasic()
	require.NoError(t, err)
}
