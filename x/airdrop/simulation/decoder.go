package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"
	"github.com/ojo-network/ojo/x/airdrop/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding airdrop type.
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.ParamsKey):
			var paramsA, paramsB types.Params
			cdc.MustUnmarshal(kvA.Value, &paramsA)
			cdc.MustUnmarshal(kvB.Value, &paramsB)
			return fmt.Sprintf("%v\n%v", paramsA, paramsB)

		case bytes.Equal(kvA.Key[:1], types.AirdropAccountKeyPrefix):
			var airdropAccountA, airdropAccountB types.AirdropAccount
			cdc.MustUnmarshal(kvA.Value, &airdropAccountA)
			cdc.MustUnmarshal(kvB.Value, &airdropAccountB)
			return fmt.Sprintf("%v\n%v", airdropAccountA, airdropAccountB)

		default:
			panic(fmt.Sprintf("invalid airdrop key prefix %X", kvA.Key[:1]))
		}
	}
}
