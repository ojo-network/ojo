package grpc

import (
	"context"
	"errors"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ojo-network/ojo/client"
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

	denomMandatoryList := []string{}
	for _, mandatoryItem := range params.MandatoryList {
		denomMandatoryList = append(denomMandatoryList, strings.ToUpper(mandatoryItem.SymbolDenom))
	}

	chainHeight, err := val1Client.NewChainHeight(ctx, zerolog.Nop())
	if err != nil {
		return err
	}

	var exchangeRates sdk.DecCoins
	for i := 0; i < 40; i++ {
		exchangeRates, err = val1Client.QueryClient.QueryExchangeRates()
		if err == nil && len(exchangeRates) == len(denomMandatoryList) {
			break
		}
		<-chainHeight.HeightChanged
	}
	// error if the loop above didn't succeed
	if err != nil {
		return err
	}
	if len(exchangeRates) != len(denomMandatoryList) {
		// TODO - update the output to display which denoms are missing https://github.com/ojo-network/ojo/issues/130
		return errors.New("couldn't fetch exchange rates matching denom accept list")
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
