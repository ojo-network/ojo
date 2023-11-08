package types

import (
	"math/big"
)

// PriceFeedData is a struct to represent the data that is relayed to other chains.
// It contains the asset name, value, resolve time, and id.
// The AssetName is an array of bytes, not a list, because lists are not
// compatible with ABI encoding.
type PriceFeedData struct {
	AssetName   [32]byte
	Value       *big.Int
	ResolveTime *big.Int
	Id          *big.Int
}

func NewPriceFeedData(
	assetName string,
	value *big.Int,
	resolveTime *big.Int,
	id *big.Int,
) PriceFeedData {
	// convert assetName to an array of bytes
	var assetNameBytes [32]byte
	copy(assetNameBytes[:], []byte(assetName))
	return PriceFeedData{
		AssetName:   assetNameBytes,
		Value:       value,
		ResolveTime: resolveTime,
		Id:          id,
	}
}
