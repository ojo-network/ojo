package abci_test

import (
	// "context"
	"testing"
	// "time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ojo-network/ojo/tests/integration"
	// "github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"

	ojoapp "github.com/ojo-network/ojo/app"
	// "github.com/ojo-network/ojo/pricefeeder"
	// "github.com/ojo-network/price-feeder/oracle"
	// "github.com/ojo-network/price-feeder/oracle/client"
	// "github.com/ojo-network/price-feeder/oracle/provider"
	// "github.com/ojo-network/price-feeder/oracle/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	ctx  sdk.Context
	app  *ojoapp.App
	keys []integration.TestValidatorKey
}

func (s *IntegrationTestSuite) SetupTest() {
	s.app, s.ctx, s.keys = integration.SetupAppWithContext(s.T())
	s.app.OracleKeeper.SetVoteThreshold(s.ctx, math.LegacyMustNewDecFromStr("0.4"))
}

func TestAbciTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

// func MockPriceFeeder(ctx context.Context) *pricefeeder.PriceFeeder {
// 	providers := make(map[types.ProviderName][]types.CurrencyPair)
// 	deviations := make(map[string]math.LegacyDec)
// 	providerEndpointsMap := make(map[types.ProviderName]provider.Endpoint)

// 	oracle := oracle.New(
// 		zerolog.Nop(),
// 		client.OracleClient{},
// 		providers,
// 		time.Second*10,
// 		deviations,
// 		providerEndpointsMap,
// 		false,
// 	)

// 	oracle.SetPrices(ctx)

// 	return &pricefeeder.PriceFeeder{
// 		Oracle: oracle,
// 	}
// }

// type MockProvider struct {}

// func (mp *MockProvider) GetTickerPrices(pairs ...types.CurrencyPair) (types.CurrencyPairTickers, error) {
// 	tickerPrices := make(types.CurrencyPairTickers, len(pairs))

// }

// func (mp *MockProvider) GetCandlePrices(pairs ...types.CurrencyPair) (types.CurrencyPairCandles, error) {

// }

// func NewMockProvider() MockProvider {
// 	ojoPair := types.CurrencyPair{
// 		Base: "OJO",
// 		Quote: "USD",
// 		Address: "OJOADDRESS",
// 	}
// 	usdcPair := types.CurrencyPair{
// 		Base: "USDC",
// 		Quote: "USD",
// 		Address: "USDCADDRESS",
// 	}
// }
