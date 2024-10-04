package types

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestPayment_TriggerUpdate(t *testing.T) {
	tests := []struct {
		name    string
		payment Payment
		rate    math.LegacyDec
		ctx     sdk.Context
		want    bool
	}{
		{
			name: "should trigger update - price deviation",
			payment: Payment{
				LastPrice: math.LegacyMustNewDecFromStr("100"),
				Deviation: math.LegacyMustNewDecFromStr("1"),
				Heartbeat: 100,
				LastBlock: 1,
			},
			rate: math.LegacyMustNewDecFromStr("101.1"),
			ctx:  sdk.Context{}.WithBlockHeight(1),
			want: true,
		},
		{
			name: "should not trigger update - no expiration",
			payment: Payment{
				LastPrice: math.LegacyMustNewDecFromStr("100"),
				Deviation: math.LegacyMustNewDecFromStr("1"),
				Heartbeat: 100,
				LastBlock: 100,
			},
			rate: math.LegacyMustNewDecFromStr("101"),
			ctx:  sdk.Context{},
			want: false,
		},
		{
			name: "should trigger update - heartbeat expired",
			payment: Payment{
				LastPrice: math.LegacyMustNewDecFromStr("100"),
				Deviation: math.LegacyMustNewDecFromStr("1"),
				Heartbeat: 100,
				LastBlock: 1,
			},
			rate: math.LegacyMustNewDecFromStr("100"),
			ctx:  sdk.Context{}.WithBlockHeight(102),
			want: true,
		},

		{
			name: "should not trigger update - deviation within threshol",
			payment: Payment{
				LastPrice: math.LegacyMustNewDecFromStr("200"),
				Deviation: math.LegacyMustNewDecFromStr("1"),
				Heartbeat: 100,
				LastBlock: 100,
			},
			rate: math.LegacyMustNewDecFromStr("202"),
			ctx:  sdk.Context{}.WithBlockHeight(101),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.payment.TriggerUpdate(tt.rate, tt.ctx)
			require.Equal(t, tt.want, got)
		})
	}
}
