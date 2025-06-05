# Oracle Module Queries

This document provides examples and explanations for querying the Ojo oracle module. The oracle module provides various query endpoints to retrieve price data, amm/accounted pool information, and other oracle-related data.

## Using the CLI

The oracle module provides several CLI commands for querying data. Here are the main query commands available:

```bash
# Query all exchange rates
ojod query oracle exchange-rates

# Query exchange rate for a specific denom
ojod query oracle exchange-rate ATOM

# Query latest prices for specific denoms
ojod query oracle latest-prices ATOM,BTC,ETH

# Query all latest prices
ojod query oracle all-latest-prices

# Query price history for an asset
ojod query oracle price-history ATOM

# Query oracle parameters
ojod query oracle params

# Query asset info for a specific denom
ojod query oracle show-asset-info ATOM

# Query all asset info
ojod query oracle list-asset-info

# Query a specific pool by ID
ojod query oracle show-pool 1

# Query all pools
ojod query oracle list-pool

# Query a specific accounted pool by ID
ojod query oracle show-accounted-pool 1

# Query all accounted pools
ojod query oracle list-accounted-pool
```

## Using gRPC

For programmatic access, you can use the gRPC client. Here are examples in Go showing how to query oracle data:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"

    sdk "github.com/cosmos/cosmos-sdk/types"
    oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

func main() {
    // Connect to your node's gRPC endpoint
    grpcConn, err := grpc.Dial(
        "localhost:9090",  // Your gRPC server address
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer grpcConn.Close()

    // Create oracle query client
    queryClient := oracletypes.NewQueryClient(grpcConn)

    // Set timeout context
    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()

    // Example 1: Query Exchange Rates
    exchangeRates, err := queryClient.ExchangeRates(ctx, &oracletypes.QueryExchangeRates{})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Exchange Rates: %v\n", exchangeRates.ExchangeRates)

    // Example 2: Query Active Exchange Rates
    activeRates, err := queryClient.ActiveExchangeRates(ctx, &oracletypes.QueryActiveExchangeRates{})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Active Rates: %v\n", activeRates.ActiveRates)

    // Example 3: Query Latest Prices
    latestPrices, err := queryClient.LatestPrices(ctx, &oracletypes.QueryLatestPricesRequest{
        Denoms: []string{"ATOM", "BTC", "ETH"},
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Latest Prices: %v\n", latestPrices)

    // Example 4: Query Price History
    priceHistory, err := queryClient.PriceHistory(ctx, &oracletypes.QueryPriceHistoryRequest{
        Asset: "ATOM",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Price History: %v\n", priceHistory)

    // Example 5: Query Oracle Parameters
    params, err := queryClient.Params(ctx, &oracletypes.QueryParams{})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Oracle Parameters: %v\n", params)

    // Example 6: Query Asset Info
    assetInfo, err := queryClient.AssetInfo(ctx, &oracletypes.QueryAssetInfoRequest{
        Denom: "ATOM",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Asset Info: %v\n", assetInfo)

    // Example 7: Query All Asset Info
    allAssetInfo, err := queryClient.AssetInfoAll(ctx, &oracletypes.QueryAssetInfoAllRequest{})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("All Asset Info: %v\n", allAssetInfo)

    // Example 8: Query a specific pool
    pool, err := queryClient.Pool(ctx, &oracletypes.QueryPoolRequest{
        PoolId: 1,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Pool: %v\n", pool)

    // Example 9: Query all pools
    allPools, err := queryClient.PoolAll(ctx, &oracletypes.QueryPoolAllRequest{})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("All Pools: %v\n", allPools)

    // Example 10: Query a specific accounted pool
    accountedPool, err := queryClient.AccountedPool(ctx, &oracletypes.QueryAccountedPoolRequest{
        PoolId: 1,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Accounted Pool: %v\n", accountedPool)

    // Example 11: Query all accounted pools
    allAccountedPools, err := queryClient.AccountedPoolAll(ctx, &oracletypes.QueryAccountedPoolAllRequest{})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("All Accounted Pools: %v\n", allAccountedPools)
}
```

## REST API Endpoints

The oracle module also exposes REST API endpoints through the gRPC-gateway. Here are the available endpoints:

```
# Get exchange rates
GET /ojo/oracle/v1/denoms/exchange_rates/{denom}

# Get active exchange rates
GET /ojo/oracle/v1/denoms/active_exchange_rates

# Get oracle parameters
GET /ojo/oracle/v1/params

# Get medians for all denoms or a specific denom
GET /ojo/historacle/v1/denoms/medians

# Get median deviations for all denoms or a specific denom
GET /ojo/historacle/v1/denoms/median_deviations

# Get latest prices for specific denoms
GET /elys-network/elys/oracle/latest_prices/{denoms}

# Get all latest prices
GET /elys-network/elys/oracle/all_latest_prices

# Get price history for an asset
GET /elys-network/elys/oracle/price_history/{asset}

# Get pool by ID
GET /elys-network/elys/oracle/pool/{pool_id}

# Get all pools
GET /elys-network/elys/oracle/pool

# Get accounted pool by ID
GET /elys-network/elys/oracle/accounted_pool/{pool_id}

# Get all accounted pools
GET /elys-network/elys/oracle/accounted_pool

# Get asset info for a specific denom
GET /elys-network/elys/oracle/asset_info/{denom}

# Get all asset info
GET /elys-network/elys/oracle/asset_info
```

## Response Types

Here are the main response types you'll receive from the queries:

### Exchange Rates Response
```go
type QueryExchangeRatesResponse struct {
    // exchange_rates defines a list of the exchange rate for all whitelisted denoms
    ExchangeRates []sdk.DecCoin
}
```

### Latest Prices Response
```go
type QueryLatestPricesResponse struct {
    Prices []types.LatestPrice
    Pagination *query.PageResponse
}

type LatestPrice struct {
    Denom       string
    LatestPrice sdk.Dec
}
```

### Price History Response
```go
type QueryPriceHistoryResponse struct {
    Prices []types.Price
    Pagination *query.PageResponse
}

type Price struct {
    Asset     string
    Price     sdk.Dec
    Timestamp uint64
}
```

### Parameters Response
```go
type QueryParamsResponse struct {
    // params defines the parameters of the module
    Params types.Params
}
```

### Asset Info Response
```go
type QueryAssetInfoResponse struct {
    AssetInfo types.AssetInfo
}

type AssetInfo struct {
    Denom   string
    Display string
    Ticker  string
    Decimal uint64
}
```

### All Asset Info Response
```go
type QueryAssetInfoAllResponse struct {
    AssetInfo []types.AssetInfo
}
```

### Pool Response
```go
type QueryPoolResponse struct {
    Pool Pool
}

type Pool struct {
    PoolId     uint64
    PoolAssets []PoolAsset
}

type PoolAsset struct {
    Token                  sdk.Coin
    Weight                 math.Int
    ExternalLiquidityRatio math.LegacyDec
}
```

### Accounted Pool Response
```go
type QueryAccountedPoolResponse struct {
    AccountedPool AccountedPool
}

type AccountedPool struct {
    PoolId      uint64
    TotalTokens []sdk.Coin
}
```

## Error Handling

When using the gRPC client, be sure to check for errors in the responses. Common error scenarios include:

- Network connectivity issues
- Invalid validator addresses
- Non-existent denoms or pool IDs
- Rate limiting
- Timeout errors

Always implement proper error handling in your applications when querying the oracle module.
