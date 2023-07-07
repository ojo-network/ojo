package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"
	ibchost "github.com/cosmos/ibc-go/v5/modules/core/24-host"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	ojoapp "github.com/ojo-network/ojo/app"
	appparams "github.com/ojo-network/ojo/app/params"

	airdroptypes "github.com/ojo-network/ojo/x/airdrop/types"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

func init() {
	simapp.GetSimulatorFlags()
}

// TestFullAppSimulation tests application fuzzing given a random seed as input.
func TestFullAppSimulation(t *testing.T) {
	config, db, dir, logger, skip, err := simapp.SetupSimulation("leveldb-app-sim", "Simulation")
	if skip {
		t.Skip("skipping application simulation")
	}

	require.NoError(t, err, "simulation setup failed")

	defer func() {
		db.Close()
		require.NoError(t, os.RemoveAll(dir))
	}()

	app := ojoapp.New(
		logger,
		db,
		nil,
		true,
		map[int64]bool{},
		ojoapp.DefaultNodeHome,
		simapp.FlagPeriodValue,
		ojoapp.MakeEncodingConfig(),
		ojoapp.EmptyAppOptions{},
		fauxMerkleModeOpt,
	)
	require.Equal(t, appparams.Name, app.Name())

	// run randomized simulation
	_, simParams, simErr := simulation.SimulateFromSeed(
		t,
		os.Stdout,
		app.BaseApp,
		initAppState(app.AppCodec(), app.StateSimulationManager),
		simtypes.RandomAccounts,
		simapp.SimulationOperations(app, app.AppCodec(), config),
		app.ModuleAccountAddrs(),
		config,
		app.AppCodec(),
	)

	// export state and simParams before the simulation error is checked
	err = simapp.CheckExportSimulation(app, config, simParams)
	require.NoError(t, err)
	require.NoError(t, simErr)

	if config.Commit {
		simapp.PrintStats(db)
	}
}

// TestAppStateDeterminism tests for application non-determinism using a PRNG
// as an input for the simulator's seed.
func TestAppStateDeterminism(t *testing.T) {
	if !simapp.FlagEnabledValue {
		t.Skip("skipping application simulation")
	}

	config := simapp.NewConfigFromFlags()
	config.InitialBlockHeight = 1
	config.ExportParamsPath = ""
	config.OnOperation = false
	config.AllInvariants = false
	config.ChainID = fmt.Sprintf("simulation-chain-%s", tmrand.NewRand().Str(6))

	numSeeds := 3
	numTimesToRunPerSeed := 5
	appHashList := make([]json.RawMessage, numTimesToRunPerSeed)

	for i := 0; i < numSeeds; i++ {
		config.Seed = rand.Int63()

		for j := 0; j < numTimesToRunPerSeed; j++ {
			var logger log.Logger
			if simapp.FlagVerboseValue {
				logger = server.ZeroLogWrapper{
					Logger: zerolog.New(os.Stderr).Level(zerolog.InfoLevel).With().Timestamp().Logger(),
				}
			} else {
				logger = server.ZeroLogWrapper{
					Logger: zerolog.Nop(),
				}
			}

			db := dbm.NewMemDB()
			app := ojoapp.New(
				logger,
				db,
				nil,
				true,
				map[int64]bool{},
				ojoapp.DefaultNodeHome,
				simapp.FlagPeriodValue,
				ojoapp.MakeEncodingConfig(),
				ojoapp.EmptyAppOptions{},
				interBlockCacheOpt(),
			)

			logger.Info(
				"running non-determinism simulation; seed %d; run: %d/%d; attempt: %d/%d\n",
				config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
			)

			_, _, err := simulation.SimulateFromSeed(
				t,
				os.Stdout,
				app.BaseApp,
				initAppState(app.AppCodec(), app.StateSimulationManager),
				simtypes.RandomAccounts,
				simapp.SimulationOperations(app, app.AppCodec(), config),
				app.ModuleAccountAddrs(),
				config,
				app.AppCodec(),
			)
			require.NoError(t, err)

			if config.Commit {
				simapp.PrintStats(db)
			}

			appHash := app.LastCommitID().Hash
			appHashList[j] = appHash

			if j != 0 {
				require.Equal(
					t,
					string(appHashList[0]),
					string(appHashList[j]),
					"non-determinism in seed %d; run: %d/%d; attempt: %d/%d\n",
					config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
				)
			}
		}
	}
}

func BenchmarkFullAppSimulation(b *testing.B) {
	config, db, dir, logger, skip, err := simapp.SetupSimulation("leveldb-app-bench-sim", "Simulation")
	if skip {
		b.Skip("skipping application simulation")
	}

	require.NoError(b, err, "simulation setup failed")

	defer func() {
		db.Close()
		require.NoError(b, os.RemoveAll(dir))
	}()

	app := ojoapp.New(
		logger,
		db,
		nil,
		true,
		map[int64]bool{},
		ojoapp.DefaultNodeHome,
		simapp.FlagPeriodValue,
		ojoapp.MakeEncodingConfig(),
		ojoapp.EmptyAppOptions{},
		interBlockCacheOpt(),
	)

	// Run randomized simulation:w
	_, simParams, simErr := simulation.SimulateFromSeed(
		b,
		os.Stdout,
		app.BaseApp,
		initAppState(app.AppCodec(), app.StateSimulationManager),
		simtypes.RandomAccounts,
		simapp.SimulationOperations(app, app.AppCodec(), config),
		app.ModuleAccountAddrs(),
		config,
		app.AppCodec(),
	)

	// export state and simParams before the simulation error is checked
	err = simapp.CheckExportSimulation(app, config, simParams)
	require.NoError(b, err)
	require.NoError(b, simErr)

	if config.Commit {
		simapp.PrintStats(db)
	}
}

func TestAppImportExport(t *testing.T) {
	db, dir, app, logger, exported, stopEarly, newDB, newDir, newApp, _ := appExportAndImport(t)
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Sprintf("%v", r)
			if !strings.Contains(err, "validator set is empty after InitGenesis") {
				panic(r)
			}
			logger.Info("Skipping simulation as all validators have been unbonded")
			logger.Info("err", err, "stacktrace", string(debug.Stack()))
		}
	}()

	defer func() {
		defer func() {
			db.Close()
			require.NoError(t, os.RemoveAll(dir))
		}()
	}()

	defer func() {
		require.NoError(t, newDB.Close())
		require.NoError(t, os.RemoveAll(newDir))
	}()

	if stopEarly {
		logger.Info("can't export or import a zero-validator genesis, exiting test...")
		return
	}

	ctxA := app.NewContext(true, tmproto.Header{Height: app.LastBlockHeight()})
	ctxB := newApp.NewContext(true, tmproto.Header{Height: app.LastBlockHeight()})

	newApp.InitChainer(ctxB, abci.RequestInitChain{
		AppStateBytes:   exported.AppState,
		ConsensusParams: exported.ConsensusParams,
	})
	newApp.StoreConsensusParams(ctxB, exported.ConsensusParams)

	logger.Info("comparing stores...\n")

	storeKeysPrefixes := []StoreKeysPrefixes{
		// The airdrop module creates vesting accounts so the number of accounts will be different
		{authtypes.ModuleName, app.GetKey(authtypes.StoreKey), newApp.GetKey(authtypes.StoreKey), [][]byte{}},
		{
			stakingtypes.ModuleName,
			app.GetKey(stakingtypes.StoreKey), newApp.GetKey(stakingtypes.StoreKey),
			[][]byte{
				stakingtypes.UnbondingQueueKey, stakingtypes.RedelegationQueueKey, stakingtypes.ValidatorQueueKey,
				stakingtypes.HistoricalInfoKey,
			},
		}, // ordering may change but it doesn't matter
		{slashingtypes.ModuleName, app.GetKey(slashingtypes.StoreKey), newApp.GetKey(slashingtypes.StoreKey), [][]byte{}},
		{minttypes.ModuleName, app.GetKey(minttypes.StoreKey), newApp.GetKey(minttypes.StoreKey), [][]byte{}},
		{distrtypes.ModuleName, app.GetKey(distrtypes.StoreKey), newApp.GetKey(distrtypes.StoreKey), [][]byte{}},
		{banktypes.ModuleName, app.GetKey(banktypes.StoreKey), newApp.GetKey(banktypes.StoreKey), [][]byte{banktypes.BalancesPrefix}},
		{paramtypes.ModuleName, app.GetKey(paramtypes.StoreKey), newApp.GetKey(paramtypes.StoreKey), [][]byte{}},
		{govtypes.ModuleName, app.GetKey(govtypes.StoreKey), newApp.GetKey(govtypes.StoreKey), [][]byte{}},
		{evidencetypes.ModuleName, app.GetKey(evidencetypes.StoreKey), newApp.GetKey(evidencetypes.StoreKey), [][]byte{}},
		{capabilitytypes.ModuleName, app.GetKey(capabilitytypes.StoreKey), newApp.GetKey(capabilitytypes.StoreKey), [][]byte{}},
		{authz.ModuleName, app.GetKey(authzkeeper.StoreKey), newApp.GetKey(authzkeeper.StoreKey), [][]byte{authzkeeper.GrantKey, authzkeeper.GrantQueuePrefix}},

		{ibchost.ModuleName, app.GetKey(ibchost.StoreKey), newApp.GetKey(ibchost.StoreKey), [][]byte{}},
		{ibctransfertypes.ModuleName, app.GetKey(ibctransfertypes.StoreKey), newApp.GetKey(ibctransfertypes.StoreKey), [][]byte{}},

		// Ojo Modules
		{oracletypes.ModuleName, app.GetKey(oracletypes.StoreKey), newApp.GetKey(oracletypes.StoreKey), [][]byte{}},
		{airdroptypes.ModuleName, app.GetKey(airdroptypes.StoreKey), newApp.GetKey(airdroptypes.StoreKey), [][]byte{}},
	}

	for _, skp := range storeKeysPrefixes {
		logger.Info("comparing ", skp.A, " and ", skp.B)
		storeA := ctxA.KVStore(skp.A)
		storeB := ctxB.KVStore(skp.B)

		failedKVAs, failedKVBs := sdk.DiffKVStores(storeA, storeB, skp.Prefixes)
		require.Equal(t, len(failedKVAs), len(failedKVBs), "unequal sets of key-values to compare")

		logger.Info(
			"using module %s StoreKey compared %d different key/value pairs between %s and %s\n",
			skp.moduleName,
			len(failedKVAs),
			skp.A,
			skp.B,
		)

		require.Equal(t, 0, len(failedKVAs), simapp.GetSimulationLog(skp.A.Name(), app.SimulationManager().StoreDecoders, failedKVAs, failedKVBs))
	}
}

func TestAppSimulationAfterImport(t *testing.T) {
	db, dir, app, logger, exported, stopEarly, newDB, newDir, newApp, config := appExportAndImport(t)
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Sprintf("%v", r)
			if !strings.Contains(err, "validator set is empty after InitGenesis") {
				panic(r)
			}
			logger.Info("Skipping simulation as all validators have been unbonded")
			logger.Info("err", err, "stacktrace", string(debug.Stack()))
		}
	}()

	defer func() {
		db.Close()
		require.NoError(t, os.RemoveAll(dir))
	}()

	defer func() {
		require.NoError(t, newDB.Close())
		require.NoError(t, os.RemoveAll(newDir))
	}()

	if stopEarly {
		logger.Info("can't export or import a zero-validator genesis, exiting test...")
		return
	}

	// importing the old app genesis into new app
	ctxB := newApp.NewContext(true, tmproto.Header{Height: app.LastBlockHeight()})
	newApp.InitChainer(ctxB, abci.RequestInitChain{
		AppStateBytes:   exported.AppState,
		ConsensusParams: exported.ConsensusParams,
	})
	newApp.StoreConsensusParams(ctxB, exported.ConsensusParams)

	_, _, err := simulation.SimulateFromSeed(
		t,
		os.Stdout,
		newApp.BaseApp,
		initAppState(newApp.AppCodec(), newApp.StateSimulationManager),
		simtypes.RandomAccounts, // Replace with own random account function if using keys other than secp256k1
		simapp.SimulationOperations(newApp, newApp.AppCodec(), config),
		newApp.ModuleAccountAddrs(),
		config,
		newApp.AppCodec(),
	)
	require.NoError(t, err)
}
