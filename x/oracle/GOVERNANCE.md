# Oracle Module Governance Operations

This document provides step-by-step examples for submitting governance proposals related to the Ojo oracle module. These examples cover parameter updates, asset management, and provider management.

## Table of Contents
- [Understanding Key Types](#understanding-key-types)
- [Updating Oracle Parameters](#updating-oracle-parameters)
- [Managing Assets](#managing-assets)
- [Managing Currency Pair Providers](#managing-currency-pair-providers)
- [Managing Currency Deviation Thresholds](#managing-currency-deviation-thresholds)
- [Canceling Parameter Updates](#canceling-parameter-updates)
- [Submitting Proposals Using JSON Files](#submitting-proposals-using-json-files)

## Understanding Key Types

### CurrencyPairProvider Type

The `CurrencyPairProvider` type defines how price data is sourced for a specific currency pair. It contains the following fields:

```go
type CurrencyPairProvider struct {
    BaseDenom               string                // The base denomination (e.g., "ATOM", "ETH")
    QuoteDenom              string                // The quote denomination (e.g., "USD", "USDT")
    BaseProxyDenom          string                // Optional proxy denomination for the base asset
    QuoteProxyDenom         string                // Optional proxy denomination for the quote asset
    PoolId                  uint64                // Optional pool ID for liquidity pool based pricing
    ExternLiquidityProvider string                // External liquidity provider identifier
    CryptoCompareExchange   string                // Optional CryptoCompare exchange identifier
    PairAddress             []PairAddressProvider // Optional addresses for on-chain DEX pairs
    Providers               []string              // List of authorized price providers
}

type PairAddressProvider struct {
    Address         string        // The contract address of the DEX pair
    AddressProvider string        // The provider type (e.g., "eth-uniswap", "eth-curve")
}
```

Key points about CurrencyPairProviders:
- Each provider represents a trusted source for price data
- Multiple providers can be specified for redundancy and accuracy
- Proxy denominations allow for indirect price calculations through intermediate assets
- External liquidity providers and pool IDs enable DEX-based pricing
- CryptoCompare integration provides additional price feed options
- PairAddress is used for on-chain DEX price sources
- Common providers include:
  - Centralized exchanges (e.g., "binance", "coinbase", "kraken")
  - DEXes (e.g., "eth-uniswap", "eth-curve", "osmosis")

### CurrencyDeviationThreshold Type

The `CurrencyDeviationThreshold` type defines the maximum allowed price deviation for a specific currency. It contains:

```go
type CurrencyDeviationThreshold struct {
    BaseDenom  string    // The denomination to monitor (e.g., "ATOM", "ETH")
    Threshold  string    // Maximum allowed deviation as a string decimal (e.g., "0.02" for 2%)
}
```

Key points about CurrencyDeviationThresholds:
- Thresholds help prevent price manipulation and ensure oracle reliability
- Values are expressed as string decimals representing percentages (e.g., "0.02" = 2%)
- Different assets may have different thresholds based on their volatility
- Exceeding the threshold may trigger alerts or affect price validity
- Typical threshold values range from 1-2% for stable assets to higher values for volatile assets
- The threshold applies to deviations from the median price across all providers

### Governance Message Types

The oracle module provides several governance message types for managing the oracle system:

#### MsgGovAddDenoms

Used to add new assets to the oracle system. Contains:
```go
type MsgGovAddDenoms struct {
    Authority                    string                           // Governance account address
    Title                       string                           // Proposal title
    Description                 string                           // Proposal description
    Height                      int64                            // Block height for the update
    DenomList                   DenomList                        // List of denominations to add
    Mandatory                   bool                             // Whether the assets should be mandatory
    RewardBand                  *cosmossdk_io_math.LegacyDec     // Optional reward band for the assets (default value of 2% will be used if not specified)
    CurrencyPairProviders       CurrencyPairProvidersList        // Price providers for the new assets
    CurrencyDeviationThresholds CurrencyDeviationThresholdList   // Deviation thresholds for the new assets
}
```

#### MsgGovUpdateParams

Used to update oracle parameters. Contains:
```go
type MsgGovUpdateParams struct {
    Authority    string           // Governance account address
    Title        string           // Proposal title
    Description  string           // Proposal description
    Plan         ParamUpdatePlan  // The parameter update plan
}
```

#### MsgGovRemoveCurrencyPairProviders

Used to remove currency pair providers from the oracle. Contains:
```go
type MsgGovRemoveCurrencyPairProviders struct {
    Authority             string                    // Governance account address
    Title                 string                    // Proposal title
    Description           string                    // Proposal description
    Height                int64                     // Block height for the update
    CurrencyPairProviders CurrencyPairProvidersList // Providers to remove
}
```

#### MsgGovRemoveCurrencyDeviationThresholds

Used to remove currency deviation thresholds. Contains:
```go
type MsgGovRemoveCurrencyDeviationThresholds struct {
    Authority    string    // Governance account address
    Title        string    // Proposal title
    Description  string    // Proposal description
    Height       int64     // Block height for the update
    Currencies   []string  // List of currency denominations to remove thresholds for
}
```

#### MsgGovCancelUpdateParamPlan

Used to cancel a scheduled parameter update. Contains:
```go
type MsgGovCancelUpdateParamPlan struct {
    Authority    string  // Governance account address
    Title        string  // Proposal title
    Description  string  // Proposal description
    Height       int64   // Height of the update to cancel
}
```

Key points about governance messages:
- All messages require governance authority for execution
- Changes are typically scheduled for a future block height
- Messages include standard proposal metadata (title, description)
- Some operations can be batched (e.g., adding multiple denoms)
- Changes go through the standard governance process with voting

## Updating Oracle Parameters

CLI example:
```bash
# Submit a parameter update proposal
ojod tx gov submit-proposal \
    --type=GovUpdateParams \
    --title="Update Oracle Parameters" \
    --description="This proposal updates the vote period and threshold" \
    --authority="ojo1..." \
    --plan.keys="VotePeriod,VoteThreshold" \
    --plan.height=1000000 \
    --plan.changes.vote_period=12 \
    --plan.changes.vote_threshold=0.5 \
    --from=<key_name> \
    --chain-id=<chain_id> \
    --gas=auto \
    --fees=1000uojo
```

## Managing Assets

### Adding New Assets

CLI example:
```bash
ojod tx gov submit-proposal \
    --type=GovAddDenoms \
    --title="Add New Assets" \
    --description="Add ATOM and BTC to the oracle" \
    --authority="ojo1..." \
    --height=1000000 \
    --denom-list='[{"base_denom":"uatom","symbol_denom":"ATOM","exponent":6},{"base_denom":"ubtc","symbol_denom":"BTC","exponent":8}]' \
    --mandatory=true \
    --reward-band=0.02 \
    --currency-pair-providers='[{"base_denom":"ATOM","providers":["binance","coinbase","kraken"]},{"base_denom":"BTC","providers":["binance","coinbase","kraken"]}]' \
    --currency-deviation-thresholds='[{"base_denom":"ATOM","threshold":"2"},{"base_denom":"BTC","threshold":"2"}]' \
    --from=<key_name> \
    --chain-id=<chain_id> \
    --gas=auto \
    --fees=1000uojo
```

## Managing Currency Pair Providers

CLI example:
```bash
ojod tx gov submit-proposal \
    --type=GovRemoveCurrencyPairProviders \
    --title="Remove Currency Pair Provider" \
    --description="Remove Binance as a provider for ATOM" \
    --authority="ojo1..." \
    --height=1000000 \
    --currency-pair-providers='[{"base_denom":"ATOM","providers":["binance"]}]' \
    --from=<key_name> \
    --chain-id=<chain_id> \
    --gas=auto \
    --fees=1000uojo
```

## Managing Currency Deviation Thresholds

CLI example:
```bash
ojod tx gov submit-proposal \
    --type=GovRemoveCurrencyDeviationThresholds \
    --title="Remove Currency Deviation Thresholds" \
    --description="Remove deviation thresholds for ATOM" \
    --authority="ojo1..." \
    --height=1000000 \
    --currencies=ATOM \
    --from=<key_name> \
    --chain-id=<chain_id> \
    --gas=auto \
    --fees=1000uojo
```

## Canceling Parameter Updates

CLI example:
```bash
ojod tx gov submit-proposal \
    --type=GovCancelUpdateParamPlan \
    --title="Cancel Parameter Update" \
    --description="Cancel the scheduled parameter update at height 1000000" \
    --authority="ojo1..." \
    --height=1000000 \
    --from=<key_name> \
    --chain-id=<chain_id> \
    --gas=auto \
    --fees=1000uojo
```

## Submitting Proposals Using JSON Files

For complex proposals, especially those involving multiple messages or detailed configurations, it's recommended to use JSON files. This approach provides better organization and reduces the chance of errors compared to command-line parameters.

### Basic Command Structure

```bash
ojod tx gov submit-proposal <proposal.json> \
    --from=<key_name> \
    --chain-id=<chain_id> \
    --gas=auto \
    --fees=1000uojo
```

### Real-World Examples

Here are some example JSON proposal files that demonstrate different oracle governance operations:

1. Adding New Assets:
```json
{
    "title": "Add ATOM with External Liquidity Configuration",
    "summary": "Adds ATOM with external liquidity configuration using Binance as provider",
    "messages": [
        {
            "@type": "/ojo.oracle.v1.MsgGovAddDenoms",
            "authority": "ojo10d07y265gmmuvt4z0w9aw880jnsr700jcz4krc",
            "title": "Add ATOM with External Liquidity",
            "description": "Adds ATOM with external liquidity configuration using Binance as provider",
            "height": 1000000,
            "denom_list": [
                {
                    "base_denom": "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2",
                    "symbol_denom": "ATOM",
                    "exponent": 6
                }
            ],
            "mandatory": true,
            "currency_pair_providers": [
                {
                    "base_denom": "ATOM",
                    "quote_denom": "USDT",
                    "base_proxy_denom": "ATOM",
                    "quote_proxy_denom": "USDC",
                    "extern_liquidity_provider": "binance",
                    "pool_id": 1,
                    "providers": [
                        "binance",
                        "okx",
                        "bitget",
                        "gate"
                    ]
                }
            ],
            "currency_deviation_thresholds": [
                {
                    "base_denom": "ATOM",
                    "threshold": "2"
                }
            ]
        }
    ],
    "metadata": "",
    "deposit": "10000000uojo"
}
```

Key components:
- Configures ATOM with both direct price feeds and external liquidity sources
- Uses Binance as the external liquidity provider with pool ID 1
- Sets up proxy denominations for price calculations (ATOM→USDC)
- Includes multiple price providers for redundancy

2. Updating Currency Pair Providers:
```json
{
    "title": "Updates cpp for RETH",
    "summary": "Updates cpp for multi eth provider with uniswap, camelot, balancer, pancake, and curve",
    "messages": [
    {
        "@type": "/ojo.oracle.v1.MsgGovUpdateParams",
        "authority": "ojo10d07y265gmmuvt4z0w9aw880jnsr700jcz4krc",
        "title": "Updates cpp for multi eth provider",
        "description": "Updates cpp for multi eth provider with uniswap, camelot, balancer, pancake, and curve",
        "plan":
            {
                "keys": [
                    "CurrencyPairProviders"
                ],
                "height": 6843550,
                "changes": {
                    "vote_period": "0",
                    "vote_threshold": "0.000000000000000000",
                    "reward_bands": [],
                    "reward_distribution_window": "0",
                    "accept_list": [],
                    "slash_fraction": "0.000000000000000000",
                    "slash_window": "0",
                    "min_valid_per_window": "0.000000000000000000",
                    "mandatory_list": [],
                    "historic_stamp_period": "0",
                    "median_stamp_period": "0",
                    "maximum_price_stamps": "0",
                    "maximum_median_stamps": "0",
                    "currency_pair_providers": [
                        {"base_denom":"USDT","quote_denom":"USD","pair_address":[],"providers":["kraken","coinbase","crypto","gate"]},
                        {"base_denom":"ATOM","quote_denom":"USDT","pair_address":[],"providers":["okx","bitget","gate"]},
                        {"base_denom":"ATOM","quote_denom":"USD","pair_address":[],"providers":["kraken"]},
                        {"base_denom":"ETH","quote_denom":"USDT","pair_address":[],"providers":["okx","bitget"]},
                        {"base_denom":"ETH","quote_denom":"USD","pair_address":[],"providers":["kraken"]},
                        {"base_denom":"RETH","quote_denom":"WETH","pair_address":[{"address":"0xa4e0faA58465A2D369aa21B3e42d43374c6F9613","address_provider":"eth-uniswap"}],"providers":["eth-uniswap"]},
                        {"base_denom":"WETH","quote_denom":"USDC","pair_address":[{"address":"0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640","address_provider":"eth-uniswap"}],"providers":["eth-uniswap"]},
                        {"base_denom":"CBETH","quote_denom":"WETH","pair_address":[{"address":"0x840deeef2f115cf50da625f7368c24af6fe74410","address_provider":"eth-uniswap"}],"providers":["eth-uniswap"]},
                        {"base_denom":"RE7LRT","quote_denom":"ETH","pair_address":[{"address":"0x4216d5900a6109bba48418b5e2Ab6cc4e61Cf477","address_provider":"eth-balancer"}],"providers":["eth-balancer"]},
                        {"base_denom":"SWBTC","quote_denom":"WBTC","pair_address":[],"providers":["eth-pancake", "eth-curve"]}
                    ]
                }
            }
        }
    ],
    "metadata": "",
    "deposit": "10000000uojo"
}
```

3. Removing Assets ([View Full Example](https://gist.github.com/rbajollari/074724ba3911318e439b73f669315caa)):
```json
{
    "title": "Removes STINJ asset",
    "summary": "Removes STINJ asset off accept list and mandatory list and removes STINJ currency providers and deviation threshold",
    "messages": [
        {
            "@type": "/ojo.oracle.v1.MsgGovUpdateParams",
            "authority": "ojo10d07y265gmmuvt4z0w9aw880jnsr700jcz4krc",
            "title": "Removes STINJ asset on ditto",
            "description": "Removes STINJ asset off accept list and mandatory list",
            "plan": {
                "keys": ["RewardBands", "AcceptList", "MandatoryList"],
                "height": 2828500,
                "changes": {
                    // Removed asset configurations
                }
            }
        },
        {
            "@type": "/ojo.oracle.v1.MsgGovRemoveCurrencyPairProviders",
            "authority": "ojo10d07y265gmmuvt4z0w9aw880jnsr700jcz4krc",
            "title": "Removes STINJ currency pair providers",
            "description": "Removes STINJ currency pair providers",
            "height": 2828500,
            "currency_pair_providers": [
                {
                    "base_denom": "STINJ",
                    "quote_denom": "INJ",
                    "pair_address": [],
                    "providers": ["astroport"]
                }
            ]
        },
        {
            "@type": "/ojo.oracle.v1.MsgGovRemoveCurrencyDeviationThresholds",
            "authority": "ojo10d07y265gmmuvt4z0w9aw880jnsr700jcz4krc",
            "title": "Removes STINJ currency deviation threshold",
            "description": "Removes STINJ currency deviation threshold",
            "height": 2828500,
            "currencies": ["STINJ"]
        }
    ],
    "metadata": "",
    "deposit": "10000000uojo"
}
```

### Key Components of JSON Proposals

1. **Message Types**: Each proposal must specify the correct message type in the `@type` field:
   - `/ojo.oracle.v1.MsgGovAddDenoms` for adding new assets
   - `/ojo.oracle.v1.MsgGovUpdateParams` for updating parameters
   - `/ojo.oracle.v1.MsgGovRemoveCurrencyPairProviders` for removing providers
   - `/ojo.oracle.v1.MsgGovRemoveCurrencyDeviationThresholds` for removing thresholds

2. **Authority**: The governance module account address that has the authority to make these changes.

3. **Height**: The block height at which the changes should take effect.

4. **Deposit**: The amount of tokens to deposit for the proposal (required for governance).

### Best Practices for JSON Proposals

1. **Validation**: Always validate your JSON file before submission:
   ```bash
   cat proposal.json | jq
   ```

2. **Multiple Messages**: When removing assets, consider grouping related operations:
   - Removing from accept/mandatory lists
   - Removing currency pair providers
   - Removing deviation thresholds

3. **Testing**: Test proposals on a testnet first, especially for complex changes.

4. **Documentation**: Include clear titles and descriptions that explain:
   - What changes are being made
   - Why the changes are necessary
   - Any potential impacts

5. **Height Planning**: Set appropriate heights that allow enough time for:
   - Proposal voting period
   - Implementation preparation
   - Validator updates

### Submitting the Proposal

```bash
# Submit the proposal
ojod tx gov submit-proposal proposal.json \
    --from=validator \
    --chain-id=ojo-1 \
    --gas=auto \
    --fees=1000uojo

# Check the proposal status
ojod query gov proposal <proposal-id>

# Vote on the proposal (after checking its content)
ojod tx gov vote <proposal-id> yes \
    --from=validator \
    --chain-id=ojo-1 \
    --gas=auto \
    --fees=1000uojo
```

## Important Notes

1. All governance proposals require the proper authority address (typically the governance module account).
2. The height parameter specifies when the changes should take effect.
3. For asset additions:
   - Make sure to specify both the base and symbol denominations correctly
   - Consider whether the asset should be mandatory
   - Set appropriate reward bands and deviation thresholds
4. When removing providers or thresholds:
   - Verify that alternatives are in place
   - Consider the impact on oracle reliability
5. Parameter updates can be scheduled for future blocks and canceled if needed
6. Always verify the proposal contents before submission
7. Follow proper governance procedures including:
   - Deposit period
   - Voting period
   - Required quorum and threshold
