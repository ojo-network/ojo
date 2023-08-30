package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCurrencyPairProvidersString(t *testing.T) {
	cpp := CurrencyPairProviders{
		BaseDenom:  "OJO",
		QuoteDenom: "USD",
		Providers: []string{
			"binance",
			"coinbase",
		},
	}
	require.Equal(t, cpp.String(), "base_denom: OJO\nquote_denom: USD\nproviders:\n    - binance\n    - coinbase\n")

	cppl := CurrencyPairProvidersList{cpp}
	require.Equal(t, cppl.String(), "base_denom: OJO\nquote_denom: USD\nproviders:\n    - binance\n    - coinbase")
}

func TestCurrencyPairProvidersEqual(t *testing.T) {
	cpp1 := CurrencyPairProviders{
		BaseDenom:  "OJO",
		QuoteDenom: "USD",
		Providers: []string{
			"binance",
			"coinbase",
		},
	}
	cpp2 := CurrencyPairProviders{
		BaseDenom:  "OJO",
		QuoteDenom: "USD",
		Providers: []string{
			"binance",
			"coinbase",
		},
	}
	cpp3 := CurrencyPairProviders{
		BaseDenom:  "OJO",
		QuoteDenom: "ATOM",
		Providers: []string{
			"binance",
			"coinbase",
		},
	}
	cpp4 := CurrencyPairProviders{
		BaseDenom:  "OJO",
		QuoteDenom: "USD",
		Providers: []string{
			"binance",
		},
	}

	require.True(t, cpp1.Equal(&cpp2))
	require.False(t, cpp1.Equal(&cpp3))
	require.False(t, cpp2.Equal(&cpp3))
	require.False(t, cpp1.Equal(&cpp4))
	require.False(t, cpp3.Equal(&cpp4))
}
