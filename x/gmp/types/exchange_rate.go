package types

import (
	"fmt"
	"math/big"
)

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
func NewPriceFeedData(
	assetName string,
	value *big.Int,
	resolveTime *big.Int,
	id *big.Int,
) (PriceFeedData, error) {
	assetSlice := []byte(assetName)
	if len(assetSlice) > 32 {
		return PriceFeedData{}, fmt.Errorf(
			"failed to parse pruning options from flags: %s", assetName,
		)
	}
	var assetArray [32]byte
	copy(assetArray[:], assetSlice)
	return PriceFeedData{
		AssetName:   assetArray,
		Value:       value,
		ResolveTime: resolveTime,
		Id:          id,
	}, nil
}
