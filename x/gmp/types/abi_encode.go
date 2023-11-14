package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// GmpEncoder is the struct we use to encode the data we want to send to the GMP.
type GmpEncoder struct {
	PriceData       []PriceData
	AssetNames      [32]byte
	ContractAddress common.Address
	CommandSelector [4]byte
	CommandParams   []byte
}

// MedianData is the struct that represents the MedianData tuple in Solidity.
type MedianData struct {
	BlockNums  []*big.Int
	Medians    []*big.Int
	Deviations []*big.Int
}

// PriceData is the struct that represents the PriceData tuple in Solidity.
type PriceData struct {
	AssetName   [32]byte
	Price       *big.Int
	ResolveTime *big.Int
	MedianData  []MedianData
}

// encoderSpec is the ABI specification for the GMP data.
var encoderSpec = abi.Arguments{
	{
		Type: priceDataType,
	},
	{
		Type: assetNameType,
	},
	{
		Type: contractAddressType,
	},
	{
		Type: commandSelectorType,
	},
	{
		Type: commandParamsType,
	},
}

// GMPEncode encodes the GMP data into a byte array.
func (g GmpEncoder) GMPEncode() ([]byte, error) {
	return encoderSpec.Pack(g.PriceData, g.AssetNames, g.ContractAddress, g.CommandSelector, g.CommandParams)
}
