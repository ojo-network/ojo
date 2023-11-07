package types

import (
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

var RateFactor = sdk.NewDec(10).Power(9)

func EncodeExchangeRate(rates []ExchangeRate) ([]byte, error) {
	abiDefinition := `[{"constant":false,"inputs":[{"name":"rates","type":"tuple[]","components":[{"name":"SymbolDenom","type":"string"},{"name":"Rate","type":"uint256"}]}],"name":"setRates","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}]`
	parsedABI, err := abi.JSON(strings.NewReader(abiDefinition))
	if err != nil {
		return nil, err
	}

	var convertedRates []interface{}
	for _, rate := range rates {
		convertedRates = append(convertedRates, struct {
			SymbolDenom string
			Rate        *big.Int
		}{
			SymbolDenom: rate.SymbolDenom,
			Rate:        rate.Rate,
		})
	}

	data, err := parsedABI.Pack("rates", convertedRates)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// DecToInt multiplies amount by rate factor to make it compatible with contracts
func DecToInt(amount sdk.Dec) *big.Int {
	return amount.Mul(RateFactor).TruncateInt().BigInt()
}
