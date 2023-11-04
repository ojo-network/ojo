package app

import (
	"encoding/json"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/ibc-go/v7/modules/apps/transfer"

	appparams "github.com/ojo-network/ojo/app/params"
)

// BankModule defines a custom wrapper around the x/bank module's AppModuleBasic
// implementation to provide custom default genesis state.
type BankModule struct {
	bank.AppModuleBasic
}

// DefaultGenesis returns custom Ojo x/bank module genesis state.
func (BankModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	OjoMetadata := banktypes.Metadata{
		Description: "The native staking token of the Ojo network.",
		Base:        appparams.BondDenom,
		Name:        appparams.DisplayDenom,
		Display:     appparams.DisplayDenom,
		Symbol:      appparams.DisplayDenom,
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    appparams.BondDenom,
				Exponent: 0,
				Aliases: []string{
					"microOjo",
				},
			},
			{
				Denom:    appparams.DisplayDenom,
				Exponent: 6,
				Aliases:  []string{},
			},
		},
	}

	genState := banktypes.DefaultGenesisState()
	genState.DenomMetadata = append(genState.DenomMetadata, OjoMetadata)

	return cdc.MustMarshalJSON(genState)
}

// StakingModule defines a custom wrapper around the x/staking module's
// AppModuleBasic implementation to provide custom default genesis state.
type StakingModule struct {
	staking.AppModuleBasic
}

// DefaultGenesis returns custom Ojo x/staking module genesis state.
func (StakingModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	p := stakingtypes.DefaultParams()
	p.BondDenom = appparams.BondDenom
	return cdc.MustMarshalJSON(&stakingtypes.GenesisState{
		Params: p,
	})
}

// CrisisModule defines a custom wrapper around the x/crisis module's
// AppModuleBasic implementation to provide custom default genesis state.
type CrisisModule struct {
	crisis.AppModuleBasic
}

// DefaultGenesis returns custom Ojo x/crisis module genesis state.
func (CrisisModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(&crisistypes.GenesisState{
		ConstantFee: sdk.NewCoin(appparams.BondDenom, sdk.NewInt(1000)),
	})
}

// MintModule defines a custom wrapper around the x/mint module's
// AppModuleBasic implementation to provide custom default genesis state.
type MintModule struct {
	mint.AppModuleBasic
}

// DefaultGenesis returns custom Ojo x/mint module genesis state.
func (MintModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genState := minttypes.DefaultGenesisState()
	genState.Params.MintDenom = appparams.BondDenom

	return cdc.MustMarshalJSON(genState)
}

// GovModule defines a custom wrapper around the x/gov module's
// AppModuleBasic implementation to provide custom default genesis state.
type GovModule struct {
	gov.AppModuleBasic
}

// DefaultGenesis returns custom Ojo x/gov module genesis state.
func (GovModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	minDeposit := sdk.NewCoins(sdk.NewCoin(appparams.BondDenom, govv1.DefaultMinDepositTokens))
	genState := govv1.DefaultGenesisState()
	genState.Params.MinDeposit = minDeposit
	genState.Params.VotingPeriod = &appparams.DefaultGovPeriod

	return cdc.MustMarshalJSON(genState)
}

// SlashingModule defines a custom wrapper around the x/slashing module's
// AppModuleBasic implementation to provide custom default genesis state.
type SlashingModule struct {
	slashing.AppModuleBasic
}

// DefaultGenesis returns custom Ojo x/slashing module genesis state.
func (SlashingModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genState := slashingtypes.DefaultGenesisState()
	genState.Params.SignedBlocksWindow = 10000
	genState.Params.DowntimeJailDuration = 24 * time.Hour

	return cdc.MustMarshalJSON(genState)
}

// IBCTransferModule defines a custom wrapper around the IBC Transfer AppModuleBasic
// so that we can use a custom Keeper.
type IBCTransferModule struct {
	transfer.AppModuleBasic
	keeper IBCTransferKeeper
}

// NewIBCTransferModule creates a new 20-transfer module
func NewIBCTransferModule(k IBCTransferKeeper) IBCTransferModule {
	return IBCTransferModule{
		keeper: k,
	}
}

// IBCAppModule is a custom wrapper around IBCModule, which
// implements the ICS26 interface for transfer given the transfer keeper.
type IBCAppModule struct {
	transfer.IBCModule
	keeper IBCTransferKeeper
}

// NewIBCAppModule creates a new IBCModule given the keeper
func NewIBCAppModule(k IBCTransferKeeper) IBCAppModule {
	return IBCAppModule{
		keeper: k,
	}
}
