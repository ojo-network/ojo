package types

import (
	"math/big"
)

type ExchangeRate struct {
	SymbolDenom string
	Rate        *big.Int
}
