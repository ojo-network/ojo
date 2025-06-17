package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCurrencyPairProvidersString(t *testing.T) {
	cpp := CurrencyPairProviders{
		BaseDenom:  "RETH",
		QuoteDenom: "WETH",
		PairAddress: []PairAddressProvider{
			{
				Address:         "address",
				AddressProvider: "eth-uniswap",
			},
		},
		Providers: []string{
			"eth-uniswap",
		},
	}
	require.Equal(
		t,
		cpp.String(),
		"base_denom: RETH\nquote_denom: WETH\nbase_proxy_denom: \"\"\nquote_proxy_denom: \"\"\npool_id: 0\nextern_liquidity_provider: \"\"\ncrypto_compare_exchange: \"\"\npair_address:\n    - address: address\n      address_provider: eth-uniswap\nproviders:\n    - eth-uniswap\n",
	)

	cppl := CurrencyPairProvidersList{cpp}
	require.Equal(
		t,
		cppl.String(),
		"base_denom: RETH\nquote_denom: WETH\nbase_proxy_denom: \"\"\nquote_proxy_denom: \"\"\npool_id: 0\nextern_liquidity_provider: \"\"\ncrypto_compare_exchange: \"\"\npair_address:\n    - address: address\n      address_provider: eth-uniswap\nproviders:\n    - eth-uniswap",
	)
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
	cpp5 := CurrencyPairProviders{
		BaseDenom:  "OJO",
		QuoteDenom: "ATOM",
		Providers: []string{
			"binance",
			"coinbase",
		},
		PairAddress: []PairAddressProvider{
			{
				Address:         "address",
				AddressProvider: "eth-uniswap",
			},
		},
	}

	require.True(t, cpp1.Equal(&cpp2))
	require.False(t, cpp1.Equal(&cpp3))
	require.False(t, cpp2.Equal(&cpp3))
	require.False(t, cpp1.Equal(&cpp4))
	require.False(t, cpp3.Equal(&cpp4))
	require.False(t, cpp4.Equal(&cpp5))
}
