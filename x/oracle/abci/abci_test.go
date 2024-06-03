package abci_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ojo-network/ojo/tests/integration"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"

	ojoapp "github.com/ojo-network/ojo/app"
	"github.com/ojo-network/ojo/pricefeeder"
	"github.com/ojo-network/price-feeder/oracle"
	"github.com/ojo-network/price-feeder/oracle/client"
	"github.com/ojo-network/price-feeder/oracle/provider"
	"github.com/ojo-network/price-feeder/oracle/types"
)

var VoteThreshold = math.LegacyMustNewDecFromStr("0.4")

type IntegrationTestSuite struct {
	suite.Suite

	ctx  sdk.Context
	app  *ojoapp.App
	keys []integration.TestValidatorKey
}

func (s *IntegrationTestSuite) SetupTest() {
	s.app, s.ctx, s.keys = integration.SetupAppWithContext(s.T())
	s.app.OracleKeeper.SetVoteThreshold(s.ctx, VoteThreshold)
}

func TestAbciTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func MockPriceFeeder() *pricefeeder.PriceFeeder {
	providers := make(map[types.ProviderName][]types.CurrencyPair)
	deviations := make(map[string]math.LegacyDec)
	providerEndpointsMap := make(map[types.ProviderName]provider.Endpoint)

	oracle := oracle.New(
		zerolog.Nop(),
		client.OracleClient{},
		providers,
		time.Second*5,
		deviations,
		providerEndpointsMap,
		false,
	)

	return &pricefeeder.PriceFeeder{
		Oracle: oracle,
	}
}
