package grpc

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ojo-network/ojo/client"
	"github.com/ojo-network/ojo/x/oracle/types"
	"github.com/rs/zerolog"
)

// MedianCheck waits for availability of all exchange rates from the denom accept list,
// records historical stamp data based on the oracle params, computes the
// median/median deviation and then compares that to the data in the
// median/median deviation gRPC query
func MedianCheck(val1Client *client.OjoClient) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	params, err := val1Client.QueryClient.QueryParams()
	if err != nil {
		return err
	}

	chainHeight, err := val1Client.NewChainHeight(ctx, zerolog.Nop())
	if err != nil {
		return err
	}

	var exchangeRates sdk.DecCoins
	var missingDenoms []string
	for i := 0; i < 40; i++ {
		exchangeRates, err = val1Client.QueryClient.QueryExchangeRates()
		missingDenoms = findMissingDenoms(exchangeRates, params.MandatoryList)
		if err == nil && len(missingDenoms) == 0 {
			break
		}
		<-chainHeight.HeightChanged
	}
	// error if the loop above didn't succeed
	if err != nil {
		return err
	}
	if len(missingDenoms) > 0 {
		return fmt.Errorf(
			"couldn't fetch exchange rates matching denom mandatory list. Missing: %s",
			strings.Join(missingDenoms, ", "),
		)
	}

	priceStore, err := listenForPrices(val1Client, params, chainHeight)
	if err != nil {
		return err
	}
	err = priceStore.checkMedians()
	if err != nil {
		return err
	}
	err = priceStore.checkMedianDeviations()
	if err != nil {
		return err
	}

	return nil
}

func findMissingDenoms(exchangeRates sdk.DecCoins, denomList types.DenomList) []string {
	missingDenoms := []string{}
	for _, denom := range denomList {
		if exchangeRates.AmountOf(denom.SymbolDenom).IsZero() {
			missingDenoms = append(missingDenoms, denom.SymbolDenom)
		}
	}
	return missingDenoms
}
