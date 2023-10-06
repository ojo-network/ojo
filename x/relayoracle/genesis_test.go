package relayoracle_test

import (
	"fmt"
	"testing"

	tmrand "github.com/cometbft/cometbft/libs/rand"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	ojoapp "github.com/ojo-network/ojo/app"
	"github.com/ojo-network/ojo/x/relayoracle"
	"github.com/ojo-network/ojo/x/relayoracle/keeper"

	"github.com/ojo-network/ojo/x/relayoracle/types"
)

type GenesisTestSuite struct {
	suite.Suite

	ctx sdk.Context
	k   keeper.Keeper
}

func (s *GenesisTestSuite) SetupTest() {
	app := ojoapp.Setup(s.T())
	s.ctx = app.BaseApp.NewContext(false, tmproto.Header{
		ChainID: fmt.Sprintf("test-chain-%s", tmrand.Str(4)),
		Height:  9,
	})

	s.k = app.RelayOracle
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

func (s *GenesisTestSuite) Test_InitGenesis() {
	genesis := types.DefaultGenesis()
	genesis.Params.PacketTimeout = 42
	genesis.PortId = "TestInit"

	// init gen state
	relayoracle.InitGenesis(s.ctx, s.k, *genesis)

	s.Require().Equal(genesis.Params, s.k.GetParams(s.ctx), "Params mismatch after InitGenesis")
	s.Require().Equal(genesis.PortId, s.k.GetPort(s.ctx), "Port ID mismatch after InitGenesis")
}

func (s *GenesisTestSuite) Test_ExportGenesis() {
	genesis := types.DefaultGenesis()
	genesis.Params.PacketTimeout = 42
	genesis.PortId = "TestExport"

	// init gen state
	relayoracle.InitGenesis(s.ctx, s.k, *genesis)

	// export gen state
	exportedGenesis := relayoracle.ExportGenesis(s.ctx, s.k)

	s.Require().Equal(genesis, exportedGenesis, "Exported genesis should match the original")
}
