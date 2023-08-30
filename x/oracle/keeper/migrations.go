package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ojo-network/ojo/x/oracle/types"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper *Keeper
}

// NewMigrator creates a Migrator.
func NewMigrator(keeper *Keeper) Migrator {
	return Migrator{keeper: keeper}
}

// MigrateValidatorSet fixes the validator set being stored as map
// causing non determinism by storing it as a list.
func (m Migrator) MigrateValidatorSet(ctx sdk.Context) {
	m.keeper.SetValidatorRewardSet(ctx)
}

// MigratePriceFeederCurrencyPairProviders adds the price feeder
// currency pair provider list.
func (m Migrator) MigratePriceFeederCurrencyPairProviders(ctx sdk.Context) {
	priceFeederCurrencyPairProviders := types.CurrencyPairProvidersList{
		types.CurrencyPairProviders{
			BaseDenom:  "USDT",
			QuoteDenom: "USD",
			Providers: []string{
				"kraken",
				"coinbase",
				"crypto",
				"gate",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "ATOM",
			QuoteDenom: "USDT",
			Providers: []string{
				"okx",
				"bitget",
				"gate",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "ATOM",
			QuoteDenom: "USD",
			Providers: []string{
				"kraken",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "ETH",
			QuoteDenom: "USDT",
			Providers: []string{
				"okx",
				"bitget",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "ETH",
			QuoteDenom: "USD",
			Providers: []string{
				"kraken",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "BTC",
			QuoteDenom: "USDT",
			Providers: []string{
				"okx",
				"gate",
				"bitget",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "BTC",
			QuoteDenom: "USD",
			Providers: []string{
				"coinbase",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "OSMO",
			QuoteDenom: "USDT",
			Providers: []string{
				"bitget",
				"gate",
				"huobi",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "OSMO",
			QuoteDenom: "ATOM",
			Providers: []string{
				"osmosisv2",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "stATOM",
			QuoteDenom: "ATOM",
			Providers: []string{
				"osmosisv2",
				"crescent",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "stOSMO",
			QuoteDenom: "OSMO",
			Providers: []string{
				"osmosisv2",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "DAI",
			QuoteDenom: "USDT",
			Providers: []string{
				"okx",
				"bitget",
				"huobi",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "DAI",
			QuoteDenom: "USD",
			Providers: []string{
				"kraken",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "JUNO",
			QuoteDenom: "USDT",
			Providers: []string{
				"bitget",
				"mexc",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "JUNO",
			QuoteDenom: "USD",
			Providers: []string{
				"kraken",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "JUNO",
			QuoteDenom: "ATOM",
			Providers: []string{
				"osmosisv2",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "stJUNO",
			QuoteDenom: "JUNO",
			Providers: []string{
				"osmosisv2",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "SCRT",
			QuoteDenom: "USD",
			Providers: []string{
				"kraken",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "SCRT",
			QuoteDenom: "USDT",
			Providers: []string{
				"mexc",
				"gate",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "WBTC",
			QuoteDenom: "USDT",
			Providers: []string{
				"okx",
				"bitget",
				"crypto",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "USDC",
			QuoteDenom: "USDT",
			Providers: []string{
				"okx",
				"bitget",
				"kraken",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "USDC",
			QuoteDenom: "USD",
			Providers: []string{
				"kraken",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "IST",
			QuoteDenom: "OSMO",
			Providers: []string{
				"osmosisv2",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "IST",
			QuoteDenom: "USDC",
			Providers: []string{
				"crescent",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "BNB",
			QuoteDenom: "USDT",
			Providers: []string{
				"mexc",
				"bitget",
				"okx",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "LUNA",
			QuoteDenom: "USDT",
			Providers: []string{
				"okx",
				"gate",
				"huobi",
				"bitget",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "DOT",
			QuoteDenom: "USD",
			Providers: []string{
				"kraken",
				"coinbase",
				"crypto",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "DOT",
			QuoteDenom: "USDT",
			Providers: []string{
				"gate",
				"bitget",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "AXL",
			QuoteDenom: "USD",
			Providers: []string{
				"coinbase",
				"crypto",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "AXL",
			QuoteDenom: "OSMO",
			Providers: []string{
				"osmosisv2",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "STARS",
			QuoteDenom: "ATOM",
			Providers: []string{
				"osmosisv2",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "STARS",
			QuoteDenom: "OSMO",
			Providers: []string{
				"osmosisv2",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "XRP",
			QuoteDenom: "USD",
			Providers: []string{
				"kraken",
				"coinbase",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "XRP",
			QuoteDenom: "USDT",
			Providers: []string{
				"gate",
				"mexc",
				"bitget",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "USK",
			QuoteDenom: "USDC",
			Providers: []string{
				"kujira",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "KUJI",
			QuoteDenom: "USDC",
			Providers: []string{
				"kujira",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "MNTA",
			QuoteDenom: "USDC",
			Providers: []string{
				"kujira",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "RETH",
			QuoteDenom: "WETH",
			Providers: []string{
				"eth-uniswap",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "WETH",
			QuoteDenom: "USDC",
			Providers: []string{
				"eth-uniswap",
			},
		},
		types.CurrencyPairProviders{
			BaseDenom:  "CBETH",
			QuoteDenom: "WETH",
			Providers: []string{
				"eth-uniswap",
			},
		},
	}
	m.keeper.SetPriceFeederCurrencyPairProvidersList(ctx, priceFeederCurrencyPairProviders)
}
