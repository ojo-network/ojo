package types

import (
	"fmt"
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

var rateFactor = sdk.NewDec(10).Power(9)

type FnName string

const (
	prices     FnName = "postPrices"
	medians    FnName = "postMedians"
	deviations FnName = "postDeviations"
)

// EncodeABI encodes the function name and parameters into ABI encoding using the above JSON.
// It can only be used with the following functions:
// - postPrices
// - postMedians
// - postDeviations
func EncodeABI(fn string, params ...interface{}) ([]byte, error) {
	switch FnName(fn) {
	case prices, medians, deviations:
		parsedABI, err := abi.JSON(strings.NewReader(abiJson))
		if err != nil {
			return nil, err
		}

		data, err := parsedABI.Pack(fn, params...)
		if err != nil {
			return nil, err
		}

		return data, nil

	default:
		return []byte{}, fmt.Errorf("invalid function name")
	}
}

// DecToInt multiplies amount by rate factor to make it compatible with contracts.
func DecToInt(amount sdk.Dec) *big.Int {
	return amount.Mul(rateFactor).TruncateInt().BigInt()
}
