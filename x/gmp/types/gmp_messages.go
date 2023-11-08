package types

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// TypeUnrecognized means coin type is unrecognized
	TypeUnrecognized = iota
	// TypeGeneralMessage is a pure message
	TypeGeneralMessage
	// TypeGeneralMessageWithToken is a general message with token
	TypeGeneralMessageWithToken
	// TypeSendToken is a direct token transfer
	TypeSendToken
)

var rateFactor = sdk.NewDec(10).Power(9)

// PriceFeedData is a struct to represent the data that is relayed to other chains.
// It contains the asset name, value, resolve time, and id.
// The AssetName is an array of bytes, not a list, because lists are not
// compatible with ABI encoding.
// Note: the ID field here is declared as "Id" because of the ABI encoding.
type PriceFeedData struct {
	AssetName   [32]byte
	Value       *big.Int
	ResolveTime *big.Int
	//nolint:stylecheck
	Id *big.Int
}

// NewPriceFeedData creates a new PriceFeedData struct.
// It must convert the assetName string to a byte array.
// This array may not exceed 32 bytes.
// TODO: Add a test for this function.
func NewPriceFeedData(
	assetName string,
	value sdk.Dec,
	resolveTime *big.Int,
	id *big.Int,
) (PriceFeedData, error) {
	assetSlice := []byte(assetName)
	if len(assetSlice) > 32 {
		return PriceFeedData{}, fmt.Errorf(
			"asset name is too long to convert to array: %s", assetName,
		)
	}
	var assetArray [32]byte
	copy(assetArray[:], assetSlice)
	return PriceFeedData{
		AssetName:   assetArray,
		Value:       decToInt(value),
		ResolveTime: resolveTime,
		Id:          id,
	}, nil
}

// DecToInt multiplies amount by rate factor to make it compatible with contracts.
func decToInt(amount sdk.Dec) *big.Int {
	return amount.Mul(rateFactor).TruncateInt().BigInt()
}
