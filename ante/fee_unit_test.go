package ante

import (
	"testing"

	evidence "cosmossdk.io/x/evidence/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestPriority(t *testing.T) {
	tcs := []struct {
		name     string
		oracle   bool
		msgs     []sdk.Msg
		priority int64
	}{
		{"empty priority 0", false, []sdk.Msg{}, 0},
		{"when oracle is set, then tx is max", true, []sdk.Msg{}, 100},
		{"evidence1", true, []sdk.Msg{&evidence.MsgSubmitEvidence{}}, 100},
		{"evidence2", false, []sdk.Msg{&evidence.MsgSubmitEvidence{}}, 90},
		{"evidence3", false, []sdk.Msg{&evidence.MsgSubmitEvidence{}, &evidence.MsgSubmitEvidence{}}, 90},
	}

	for _, tc := range tcs {
		p := getTxPriority(tc.oracle, tc.msgs)
		assert.Equal(t, tc.priority, p, tc.name)
	}
}
