package types

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

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
		parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
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
