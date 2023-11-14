package types

import (
	fmt "fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
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

// GmpEncoder is the struct we use to encode the data we want to send to the GMP.
type GmpEncoder struct {
	PriceData       []PriceData
	AssetNames      [][32]byte
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
		Type: assetNamesType,
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

func NewGMPEncoder(
	priceData []PriceData,
	assetName []string,
	contractAddress common.Address,
	commandSelector [4]byte,
	commandParams []byte,
) GmpEncoder {
	return GmpEncoder{
		PriceData:       priceData,
		AssetNames:      namesToBytes(assetName),
		ContractAddress: contractAddress,
		CommandSelector: commandSelector,
		CommandParams:   commandParams,
	}
}

func nameToBytes32(name string) [32]byte {
	var nameBytes [32]byte
	copy(nameBytes[:], []byte(name))
	return nameBytes
}

func namesToBytes(assetNames []string) [][32]byte {
	assetNamesBytes := make([][32]byte, len(assetNames))
	for i, name := range assetNames {
		assetNamesBytes[i] = nameToBytes32(name)
	}
	return assetNamesBytes
}

func NewPriceData(
	assetName string,
	price sdk.Dec,
	resolveTime *big.Int,
	medianData []MedianData,
) (PriceData, error) {
	assetSlice := []byte(assetName)
	if len(assetSlice) > 32 {
		return PriceData{}, fmt.Errorf(
			"asset name is too long to convert to array: %s", assetName,
		)
	}
	var assetArray [32]byte
	copy(assetArray[:], assetSlice)
	return PriceData{
		AssetName:   assetArray,
		Price:       decToInt(price),
		ResolveTime: resolveTime,
		MedianData:  medianData,
	}, nil
}

// DecToInt multiplies amount by rate factor to make it compatible with contracts.
func decToInt(amount sdk.Dec) *big.Int {
	return amount.Mul(rateFactor).TruncateInt().BigInt()
}

var rateFactor = sdk.NewDec(10).Power(9)

// NewMediansSlice creates a slice of MedianData from slices of medians and deviations.
func NewMediansSlice(medians oracletypes.PriceStamps, deviations oracletypes.PriceStamps) ([]MedianData, error) {
	if len(medians) != len(deviations) {
		return nil, fmt.Errorf("length of medians and deviations must be equal")
	}
	// First, sort them so we'll be able to find matching blockNums.
	sortedMedians := *medians.Sort()
	sortedDeviations := *deviations.Sort()

	// Then, create the MedianData slice.
	medianData := make([]MedianData, 0, len(medians))
	for i, median := range sortedMedians {
		// If the median and deviation are not from the same block, skip.
		if median.BlockNum != sortedDeviations[i].BlockNum {
			continue
		}
		medianData = append(medianData, MedianData{
			BlockNums:  []*big.Int{big.NewInt(int64(median.BlockNum))},
			Medians:    []*big.Int{decToInt(median.ExchangeRate.Amount)},
			Deviations: []*big.Int{decToInt(sortedDeviations[i].ExchangeRate.Amount)},
		})
	}
	return medianData, nil
}