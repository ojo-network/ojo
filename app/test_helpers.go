package app

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	pruningtypes "cosmossdk.io/store/pruning/types"
	//"cosmossdk.io/x/tx/signing"
	abci "github.com/cometbft/cometbft/abci/types"
	cmttypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	//"github.com/cosmos/cosmos-sdk/codec/address"
	//"github.com/cosmos/cosmos-sdk/codec/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	//"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	//authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	//"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"

	"github.com/ojo-network/ojo/app/params"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

type EmptyAppOptions struct{}

func (EmptyAppOptions) Get(_ string) interface{} { return nil }

func Setup(t *testing.T) *App {
	t.Helper()

	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err)

	// create validator set with single validator
	validator := cmttypes.NewValidator(pubKey, 1)
	valSet := cmttypes.NewValidatorSet([]*cmttypes.Validator{validator})

	// generate genesis account
	senderPrivKey := secp256k1.GenPrivKey()
	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(params.BondDenom, math.NewInt(10000000000000000))),
	}

	app := SetupWithGenesisValSet(t, valSet, []authtypes.GenesisAccount{acc}, balance)
	return app
}

// SetupWithGenesisValSet initializes a new app with a validator set and genesis accounts
// that also act as delegators. For simplicity, each validator is bonded with a delegation
// of one consensus engine unit in the default token of the simapp from first genesis
// account. A Nop logger is set in app.
func SetupWithGenesisValSet(
	t *testing.T,
	valSet *cmttypes.ValidatorSet,
	genAccs []authtypes.GenesisAccount,
	balances ...banktypes.Balance,
) *App {
	t.Helper()

	app, genesisState := setup(true, 5)
	genesisState, err := GenesisStateWithValSet(app.AppCodec(), genesisState, valSet, genAccs, balances...)
	require.NoError(t, err)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	require.NoError(t, err)

	// init chain will set the validator set and initialize the genesis accounts
	_, err = app.InitChain(
		&abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: simtestutil.DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)
	require.NoError(t, err)

	_, err = app.FinalizeBlock(&abci.RequestFinalizeBlock{
		Height:             app.LastBlockHeight() + 1,
		Hash:               app.LastCommitID().Hash,
		NextValidatorsHash: valSet.Hash(),
	})
	require.NoError(t, err)

	return app
}

// GenesisStateWithValSet returns a new genesis state with the validator set
func GenesisStateWithValSet(codec codec.Codec, genesisState map[string]json.RawMessage,
	valSet *cmttypes.ValidatorSet, genAccs []authtypes.GenesisAccount,
	balances ...banktypes.Balance,
) (map[string]json.RawMessage, error) {
	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = codec.MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := sdk.DefaultPowerReduction

	for _, val := range valSet.Validators {
		pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey)
		if err != nil {
			return nil, fmt.Errorf("failed to convert pubkey: %w", err)
		}

		pkAny, err := codectypes.NewAnyWithValue(pk)
		if err != nil {
			return nil, fmt.Errorf("failed to create new any: %w", err)
		}

		validator := stakingtypes.Validator{
			OperatorAddress:   sdk.ValAddress(val.Address).String(),
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingtypes.Bonded,
			Tokens:            bondAmt,
			DelegatorShares:   math.LegacyOneDec(),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec()),
			MinSelfDelegation: math.ZeroInt(),
		}
		validators = append(validators, validator)
		newDel := stakingtypes.NewDelegation(genAccs[0].GetAddress().String(), sdk.ValAddress(val.Address).String(), math.LegacyOneDec())
		delegations = append(delegations, newDel)

	}

	defaultStParams := stakingtypes.DefaultParams()
	stParams := stakingtypes.NewParams(
		defaultStParams.UnbondingTime,
		defaultStParams.MaxValidators,
		defaultStParams.MaxEntries,
		defaultStParams.HistoricalEntries,
		params.BondDenom,
		defaultStParams.MinCommissionRate,
	)
	// set validators and delegations
	stakingGenesis := stakingtypes.NewGenesisState(stParams, validators, delegations)
	genesisState[stakingtypes.ModuleName] = codec.MustMarshalJSON(stakingGenesis)

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		// add genesis acc tokens to total supply
		totalSupply = totalSupply.Add(b.Coins...)
	}

	for range delegations {
		// add delegated tokens to total supply
		totalSupply = totalSupply.Add(sdk.NewCoin(params.BondDenom, bondAmt))
	}

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(params.BondDenom, bondAmt)},
	})

	// update total supply
	bankGenesis := banktypes.NewGenesisState(
		banktypes.DefaultGenesisState().Params,
		balances,
		totalSupply,
		[]banktypes.Metadata{},
		[]banktypes.SendEnabled{},
	)
	genesisState[banktypes.ModuleName] = codec.MustMarshalJSON(bankGenesis)

	return genesisState, nil
}

func setup(withGenesis bool, invCheckPeriod uint) (*App, GenesisState) {
	db := dbm.NewMemDB()
	appOptions := make(simtestutil.AppOptionsMap, 0)

	app := New(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		DefaultNodeHome,
		invCheckPeriod,
		appOptions,
	)
	if withGenesis {
		return app, app.DefaultGenesis()
	}

	return app, GenesisState{}
}

// IntegrationTestNetworkConfig returns a networking configuration used for
// integration tests using the SDK's in-process network test suite.
func IntegrationTestNetworkConfig() network.Config {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(params.AccountAddressPrefix, params.AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(params.ValidatorAddressPrefix, params.ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(params.ConsNodeAddressPrefix, params.ConsNodePubKeyPrefix)
	cfg := network.DefaultConfig(NewTestNetworkFixture)
	app, genState := setup(true, 100)

	var oracleGenState oracletypes.GenesisState
	if err := app.appCodec.UnmarshalJSON(genState[oracletypes.ModuleName], &oracleGenState); err != nil {
		panic(err)
	}

	// Set mock exchange rates and a large enough vote period such that we won't
	// execute ballot voting and thus clear out previous exchange rates, since we
	// are not running a price-feeder.
	oracleGenState.Params.VotePeriod = 1000
	oracleGenState.ExchangeRates = append(oracleGenState.ExchangeRates, sdk.NewDecCoinFromDec(
		params.DisplayDenom, math.LegacyMustNewDecFromStr("34.21"),
	))

	bz, err := app.appCodec.MarshalJSON(&oracleGenState)
	if err != nil {
		panic(err)
	}
	genState[oracletypes.ModuleName] = bz

	var govGenState govv1.GenesisState
	if err := app.appCodec.UnmarshalJSON(genState[govtypes.ModuleName], &govGenState); err != nil {
		panic(err)
	}

	votingPeriod := time.Minute
	govGenState.Params.VotingPeriod = &votingPeriod

	bz, err = app.appCodec.MarshalJSON(&govGenState)
	if err != nil {
		panic(err)
	}
	genState[govtypes.ModuleName] = bz

	stakingGenesis := stakingtypes.GetGenesisStateFromAppState(app.appCodec, genState)
	stakingGenesis.Params.BondDenom = params.BondDenom
	genState[stakingtypes.ModuleName] = app.appCodec.MustMarshalJSON(stakingGenesis)

	cfg.Codec = app.appCodec
	cfg.TxConfig = app.txConfig
	cfg.LegacyAmino = app.legacyAmino
	cfg.InterfaceRegistry = app.interfaceRegistry
	cfg.GenesisState = genState
	cfg.MinGasPrices = params.ProtocolMinGasPrice.String()
	cfg.BondDenom = params.BondDenom

	return cfg
}

func NewTestNetworkFixture() network.TestFixture {
	dir, err := os.MkdirTemp("", "simapp")
	if err != nil {
		panic(fmt.Sprintf("failed creating temporary directory: %v", err))
	}
	defer os.RemoveAll(dir)

	app, genState := setup(true, 100)

	appCtr := func(val network.ValidatorI) servertypes.Application {
		return New(
			val.GetCtx().Logger,
			dbm.NewMemDB(),
			nil,
			true,
			make(map[int64]bool),
			val.GetCtx().Config.RootDir,
			0,
			EmptyAppOptions{},
			baseapp.SetPruning(pruningtypes.NewPruningOptionsFromString(val.GetAppConfig().Pruning)),
			baseapp.SetMinGasPrices(val.GetAppConfig().MinGasPrices),
			baseapp.SetChainID(val.GetCtx().Viper.GetString(flags.FlagChainID)),
		)
	}

	return network.TestFixture{
		AppConstructor: appCtr,
		GenesisState:   genState,
		EncodingConfig: testutil.TestEncodingConfig{
			InterfaceRegistry: app.InterfaceRegistry(),
			Codec:             app.AppCodec(),
			TxConfig:          app.txConfig,
			Amino:             app.LegacyAmino(),
		},
	}
}
