package main

import (
	"github.com/umee-network/umee-infra/infra/pulumi/common/jsonmutator"
	netconfig "github.com/umee-network/umee-infra/lib/umeednetworkconfigurator"
)

type Network struct {
	ChainID                 string
	NumNodes                int
	LocalCosmosBinaryPath   string
	CosmosHomeFolderName    string
	NodeConfig              NodeConfig
	NodeGenesisAccounts     []netconfig.NodeGenesisAccountConfig
	GenesisAccounts         []netconfig.GenesisAccount
	NetworkGenesisMutations NetworkGenesisMutations
}

type NetworkGenesisMutations map[string]interface{}

func (mutations NetworkGenesisMutations) MutateGenesis(originalGenesis string) (string, error) {
	g := jsonmutator.NewJSONMutator(originalGenesis)

	for key, val := range mutations {
		g.Set(key, val)
	}

	return g.Out()
}

type NodeLocation struct {
	Region string
	Zone   string
}

type NodeConfig struct {
	MachineType string
	DiskType    string
	DiskSizeGB  int
	Locations   []NodeLocation
}

type NodeSecretConfig struct{}

type Disk struct {
	Name       string
	MountPoint string
	Size       int
	Type       string
}
