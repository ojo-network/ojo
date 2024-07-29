package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ojo-network/ojo/app"
	appparams "github.com/ojo-network/ojo/app/params"
	"github.com/ojo-network/ojo/client"

	dbm "github.com/cometbft/cometbft-db"
	tmconfig "github.com/cometbft/cometbft/config"
	tmjson "github.com/cometbft/cometbft/libs/json"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/server"
	srvconfig "github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/spf13/viper"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	airdroptypes "github.com/ojo-network/ojo/x/airdrop/types"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

const (
	ojoContainerRepo  = "ojo"
	ojoP2pPort        = "26656"
	ojoTmrpcPort      = "26657"
	ojoGrpcPort       = "9090"
	ojoMaxStartupTime = 40 // seconds

	priceFeederContainerRepo  = "ghcr.io/ojo-network/price-feeder-ojo-v0.1.11"
	priceFeederServerPort     = "7171/tcp"
	priceFeederMaxStartupTime = 20 // seconds

	initBalanceStr = "510000000000" + appparams.BondDenom
)

type Orchestrator struct {
	tmpDirs             []string
	chain               *chain
	dkrPool             *dockertest.Pool
	dkrNet              *dockertest.Network
	priceFeederResource *dockertest.Resource
	OjoClient           *client.OjoClient // signs tx with val[0]
	AirdropClient       *client.OjoClient // signs tx with account[0]
}

// SetupSuite initializes and runs all the resources needed for the
// e2e integration test suite
// 1. creates a temporary directory for the chain data and initializes keys for all validators
// 2. initializes the genesis files for all validators
// 3. starts up each validator in its own docker container with the 3rd validator holding the majority of the stake
// 4. initializes the ojo client used to send transactions and queries to the validators
// 5. delegates voting power from the majority share validator to another one for the price feeder
// 6. starts up the price feeder in its own docker container
func (o *Orchestrator) InitResources(t *testing.T) {
	t.Log("setting up e2e integration test suite...")
	appparams.SetAddressPrefixes()

	db := dbm.NewMemDB()
	app := app.New(
		nil,
		db,
		nil,
		true,
		map[int64]bool{},
		"",
		0,
		app.EmptyAppOptions{},
	)
	encodingConfig = testutil.TestEncodingConfig{
		InterfaceRegistry: app.InterfaceRegistry(),
		Codec:             app.AppCodec(),
		TxConfig:          app.GetTxConfig(),
		Amino:             app.LegacyAmino(),
	}

	// codec
	cdc := encodingConfig.Codec

	var err error
	o.chain, err = newChain(cdc)
	require.NoError(t, err)

	t.Logf("starting e2e infrastructure; chain-id: %s; datadir: %s", o.chain.id, o.chain.dataDir)

	o.dkrPool, err = dockertest.NewPool("")
	require.NoError(t, err)

	o.dkrNet, err = o.dkrPool.CreateNetwork(fmt.Sprintf("%s-testnet", o.chain.id))
	require.NoError(t, err)

	o.initNodes(t)
	o.initUserAccounts(t)
	o.initGenesis(t)
	o.initValidatorConfigs(t)
	o.runValidators(t)
	o.initOjoClient(t)
	o.initAirdropClient(t)
	o.delegatePriceFeederVoting(t)
	o.runPriceFeeder(t)
}

func (o *Orchestrator) TearDownResources(t *testing.T) {
	t.Log("tearing down e2e integration test suite...")

	require.NoError(t, o.dkrPool.Purge(o.priceFeederResource))

	for _, val := range o.chain.validators {
		require.NoError(t, o.dkrPool.Purge(val.dockerResource))
	}

	require.NoError(t, o.dkrPool.RemoveNetwork(o.dkrNet))

	os.RemoveAll(o.chain.dataDir)
	for _, td := range o.tmpDirs {
		os.RemoveAll(td)
	}
}

func (o *Orchestrator) initNodes(t *testing.T) {
	require.NoError(t, o.chain.createAndInitValidators(2))

	// initialize a genesis file for the first validator
	val0ConfigDir := o.chain.validators[0].configDir()
	for _, val := range o.chain.validators {
		valAddr, err := val.keyInfo.GetAddress()
		require.NoError(t, err)
		require.NoError(t,
			addGenesisAccount(o.chain.cdc, val0ConfigDir, "", initBalanceStr, valAddr),
		)
	}

	// copy the genesis file to the remaining validators
	for _, val := range o.chain.validators[1:] {
		_, err := copyFile(
			filepath.Join(val0ConfigDir, "config", "genesis.json"),
			filepath.Join(val.configDir(), "config", "genesis.json"),
		)
		require.NoError(t, err)
	}
}

func (o *Orchestrator) initUserAccounts(t *testing.T) {
	err := o.chain.createAccounts(1)
	require.NoError(t, err)
}

func (o *Orchestrator) initGenesis(t *testing.T) {
	serverCtx := server.NewDefaultContext()
	config := serverCtx.Config

	config.SetRoot(o.chain.validators[0].configDir())
	config.Moniker = o.chain.validators[0].moniker

	genFilePath := config.GenesisFile()
	t.Log("starting e2e infrastructure; validator_0 config:", genFilePath)
	appGenState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFilePath)
	require.NoError(t, err)

	// Oracle
	var oracleGenState oracletypes.GenesisState
	require.NoError(t, o.chain.cdc.UnmarshalJSON(appGenState[oracletypes.ModuleName], &oracleGenState))

	oracleGenState.Params.HistoricStampPeriod = 3
	oracleGenState.Params.MaximumPriceStamps = 4
	oracleGenState.Params.MedianStampPeriod = 12
	oracleGenState.Params.MaximumMedianStamps = 2
	oracleGenState.Params.AcceptList = oracleAcceptList
	oracleGenState.Params.MandatoryList = oracleMandatoryList
	oracleGenState.Params.RewardBands = oracleRewardBands

	bz, err := o.chain.cdc.MarshalJSON(&oracleGenState)
	require.NoError(t, err)
	appGenState[oracletypes.ModuleName] = bz

	// Airdrop
	var airdropGenState airdroptypes.GenesisState
	require.NoError(t, o.chain.cdc.UnmarshalJSON(appGenState[airdroptypes.ModuleName], &airdropGenState))

	// Use the first and only account as the airdrop origin
	airdropOriginAddress, err := o.chain.accounts[0].KeyInfo.GetAddress()
	require.NoError(t, err)
	airdropGenState.AirdropAccounts = []*airdroptypes.AirdropAccount{
		{
			OriginAddress:  airdropOriginAddress.String(),
			OriginAmount:   100000000000,
			VestingEndTime: time.Now().Add(24 * time.Hour).Unix(),
			State:          airdroptypes.AirdropAccount_STATE_CREATED,
		},
	}

	bz, err = o.chain.cdc.MarshalJSON(&airdropGenState)
	require.NoError(t, err)
	appGenState[airdroptypes.ModuleName] = bz

	// Gov
	var govGenState govtypesv1.GenesisState
	require.NoError(t, o.chain.cdc.UnmarshalJSON(appGenState[govtypes.ModuleName], &govGenState))

	votingPeroid := 5 * time.Second
	govGenState.Params.VotingPeriod = &votingPeroid

	bz, err = o.chain.cdc.MarshalJSON(&govGenState)
	require.NoError(t, err)
	appGenState[govtypes.ModuleName] = bz

	// Genesis Txs
	var genUtilGenState genutiltypes.GenesisState
	require.NoError(t, o.chain.cdc.UnmarshalJSON(appGenState[genutiltypes.ModuleName], &genUtilGenState))

	genTxs := make([]json.RawMessage, len(o.chain.validators))
	for i, val := range o.chain.validators {
		var createValmsg sdk.Msg
		if i == 2 {
			createValmsg, err = val.buildCreateValidatorMsg(majorityValidatorStake)
		} else {
			createValmsg, err = val.buildCreateValidatorMsg(minorityValidatorStake)
		}
		require.NoError(t, err)

		signedTx, err := val.signMsg(o.chain.cdc, createValmsg)
		require.NoError(t, err)

		txRaw, err := o.chain.cdc.MarshalJSON(signedTx)
		require.NoError(t, err)

		genTxs[i] = txRaw
	}

	genUtilGenState.GenTxs = genTxs

	bz, err = o.chain.cdc.MarshalJSON(&genUtilGenState)
	require.NoError(t, err)
	appGenState[genutiltypes.ModuleName] = bz

	bz, err = json.MarshalIndent(appGenState, "", "  ")
	require.NoError(t, err)

	genDoc.AppState = bz

	bz, err = tmjson.MarshalIndent(genDoc, "", "  ")
	require.NoError(t, err)

	// write the updated genesis file to each validator
	for _, val := range o.chain.validators {
		err := writeFile(filepath.Join(val.configDir(), "config", "genesis.json"), bz)
		require.NoError(t, err)
	}
}

func (o *Orchestrator) initValidatorConfigs(t *testing.T) {
	for i, val := range o.chain.validators {
		tmCfgPath := filepath.Join(val.configDir(), "config", "config.toml")

		vpr := viper.New()
		vpr.SetConfigFile(tmCfgPath)
		require.NoError(t, vpr.ReadInConfig())

		valConfig := tmconfig.DefaultConfig()
		require.NoError(t, vpr.Unmarshal(valConfig))

		valConfig.P2P.ListenAddress = fmt.Sprintf("tcp://0.0.0.0:%s", ojoP2pPort)
		valConfig.P2P.AddrBookStrict = false
		valConfig.P2P.ExternalAddress = fmt.Sprintf("%s:%s", val.instanceName(), ojoP2pPort)
		valConfig.RPC.ListenAddress = fmt.Sprintf("tcp://0.0.0.0:%s", ojoTmrpcPort)
		valConfig.StateSync.Enable = false
		valConfig.LogLevel = "info"

		var peers []string

		for j := 0; j < len(o.chain.validators); j++ {
			if i == j {
				continue
			}

			peer := o.chain.validators[j]
			peerID := fmt.Sprintf("%s@%s%d:26656", peer.nodeKey.ID(), peer.moniker, j)
			peers = append(peers, peerID)
		}

		valConfig.P2P.PersistentPeers = strings.Join(peers, ",")

		tmconfig.WriteConfigFile(tmCfgPath, valConfig)

		// set application configuration
		appCfgPath := filepath.Join(val.configDir(), "config", "app.toml")

		appConfig := srvconfig.DefaultConfig()
		appConfig.API.Enable = true
		appConfig.API.Address = "tcp://0.0.0.0:1317"
		appConfig.MinGasPrices = minGasPrice
		appConfig.GRPC.Address = "0.0.0.0:9090"

		srvconfig.WriteConfigFile(appCfgPath, appConfig)
	}
}

func (o *Orchestrator) runValidators(t *testing.T) {
	t.Log("starting ojo validator containers...")

	proposalsDirectory, err := proposalsDirectory()
	require.NoError(t, err)

	for _, val := range o.chain.validators {
		runOpts := &dockertest.RunOptions{
			Name:      val.instanceName(),
			NetworkID: o.dkrNet.Network.ID,
			Mounts: []string{
				fmt.Sprintf("%s/:/root/.ojo", val.configDir()),
				fmt.Sprintf("%s/:/root/proposals", proposalsDirectory),
			},
			Repository: ojoContainerRepo,
		}

		// expose the first validator
		if val.index == 0 {
			runOpts.PortBindings = map[docker.Port][]docker.PortBinding{
				"1317/tcp":  {{HostIP: "", HostPort: "1317"}},
				"9090/tcp":  {{HostIP: "", HostPort: "9090"}},
				"26656/tcp": {{HostIP: "", HostPort: "26656"}},
				"26657/tcp": {{HostIP: "", HostPort: "26657"}},
			}
		}

		resource, err := o.dkrPool.RunWithOptions(runOpts, noRestart)
		require.NoError(t, err)

		val.dockerResource = resource
		t.Logf("started ojo validator container: %s", resource.Container.ID)
	}

	rpcURL := fmt.Sprintf("tcp://localhost:%s", ojoTmrpcPort)
	rpcClient, err := rpchttp.New(rpcURL, "/websocket")
	require.NoError(t, err)

	checkHealth := func() bool {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		status, err := rpcClient.Status(ctx)
		if err != nil {
			return false
		}

		// let the node produce a few blocks
		if status.SyncInfo.CatchingUp || status.SyncInfo.LatestBlockHeight < 3 {
			return false
		}

		return true
	}

	isHealthy := false
	for i := 0; i < ojoMaxStartupTime; i++ {
		isHealthy = checkHealth()
		if isHealthy {
			break
		}
		time.Sleep(time.Second)
	}

	if !isHealthy {
		err := o.outputLogs(o.chain.validators[0].dockerResource)
		if err != nil {
			t.Log("Error retrieving Ojo node logs", err)
		}
		t.Fatal("Ojo node failed to produce blocks")
	}
}

func (o *Orchestrator) delegatePriceFeederVoting(t *testing.T) {
	delegateAddr, err := o.chain.validators[0].keyInfo.GetAddress()
	require.NoError(t, err)
	_, err = o.OjoClient.TxClient.TxDelegateFeedConsent(delegateAddr)
	require.NoError(t, err)
}

func (o *Orchestrator) runPriceFeeder(t *testing.T) {
	t.Log("starting price-feeder container...")

	votingVal := o.chain.validators[1]
	votingValAddr, err := votingVal.keyInfo.GetAddress()
	require.NoError(t, err)

	delegateVal := o.chain.validators[0]
	delegateValAddr, err := delegateVal.keyInfo.GetAddress()
	require.NoError(t, err)

	grpcEndpoint := fmt.Sprintf("tcp://%s:%s", delegateVal.instanceName(), ojoGrpcPort)
	tmrpcEndpoint := fmt.Sprintf("http://%s:%s", delegateVal.instanceName(), ojoTmrpcPort)

	o.priceFeederResource, err = o.dkrPool.RunWithOptions(
		&dockertest.RunOptions{
			Name:       "price-feeder",
			NetworkID:  o.dkrNet.Network.ID,
			Repository: priceFeederContainerRepo,
			Mounts: []string{
				fmt.Sprintf("%s/:/root/.ojo", delegateVal.configDir()),
			},
			PortBindings: map[docker.Port][]docker.PortBinding{
				"7171/tcp": {{HostIP: "", HostPort: "7171"}},
			},
			Env: []string{
				fmt.Sprintf("PRICE_FEEDER_PASS=%s", keyringPassphrase),
				fmt.Sprintf("ACCOUNT_ADDRESS=%s", delegateValAddr),
				fmt.Sprintf("ACCOUNT_VALIDATOR=%s", sdk.ValAddress(votingValAddr)),
				fmt.Sprintf("KEYRING_DIR=%s", "/root/.ojo"),
				fmt.Sprintf("ACCOUNT_CHAIN_ID=%s", o.chain.id),
				fmt.Sprintf("RPC_GRPC_ENDPOINT=%s", grpcEndpoint),
				fmt.Sprintf("RPC_TMRPC_ENDPOINT=%s", tmrpcEndpoint),
			},
			Cmd: []string{"--skip-provider-check", "--config-currency-providers", "--log-level=debug"},
		},
		noRestart,
	)
	require.NoError(t, err)

	endpoint := fmt.Sprintf("http://%s/api/v1/prices", o.priceFeederResource.GetHostPort(priceFeederServerPort))

	checkHealth := func() bool {
		resp, err := http.Get(endpoint)
		if err != nil {
			t.Log("Price feeder endpoint not available", err, endpoint)
			return false
		}

		defer resp.Body.Close()

		bz, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Log("Can't get price feeder response", err)
			return false
		}

		var respBody map[string]interface{}
		if err := json.Unmarshal(bz, &respBody); err != nil {
			t.Log("Can't unmarshal price feed", err)
			return false
		}

		prices, ok := respBody["prices"].(map[string]interface{})
		if !ok {
			t.Log("price feeder: no prices")
			return false
		}

		return len(prices) > 0
	}

	isHealthy := false
	for i := 0; i < priceFeederMaxStartupTime; i++ {
		isHealthy = checkHealth()
		if isHealthy {
			break
		}
		time.Sleep(time.Second)
	}

	if !isHealthy {
		err := o.outputLogs(o.priceFeederResource)
		if err != nil {
			t.Log("Error retrieving price feeder logs", err)
		}
		t.Fatal("price-feeder not healthy")
	}

	t.Logf("started price-feeder container: %s", o.priceFeederResource.Container.ID)
}

func (o *Orchestrator) initOjoClient(t *testing.T) {
	var err error
	o.OjoClient, err = client.NewOjoClient(
		o.chain.id,
		fmt.Sprintf("tcp://localhost:%s", ojoTmrpcPort),
		fmt.Sprintf("tcp://localhost:%s", ojoGrpcPort),
		"val1",
		o.chain.validators[1].mnemonic,
	)
	require.NoError(t, err)
}

func (o *Orchestrator) initAirdropClient(t *testing.T) {
	var err error
	o.AirdropClient, err = client.NewOjoClient(
		o.chain.id,
		fmt.Sprintf("tcp://localhost:%s", ojoTmrpcPort),
		fmt.Sprintf("tcp://localhost:%s", ojoGrpcPort),
		"val1",
		o.chain.accounts[0].Mnemonic,
	)
	require.NoError(t, err)
}

func (o *Orchestrator) outputLogs(resource *dockertest.Resource) error {
	return o.dkrPool.Client.Logs(docker.LogsOptions{
		Container:    resource.Container.ID,
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
		Stdout:       true,
		Stderr:       true,
		Tail:         "false",
	})
}

func noRestart(config *docker.HostConfig) {
	config.RestartPolicy = docker.RestartPolicy{
		Name: "no",
	}
}

func proposalsDirectory() (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	adjacentDirPath := filepath.Join(workingDir, "proposals")
	absoluteAdjacentDirPath, err := filepath.Abs(adjacentDirPath)
	if err != nil {
		return "", err
	}

	return absoluteAdjacentDirPath, nil
}
