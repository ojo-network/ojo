package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/ojo-network/ojo/x/airdrop/types"
)

const (
	expiryBlockKey           = "expiry_block"
	delegationRequirementKey = "delegation_requirement"
	airdropFactorKey         = "airdrop_factor"

	numAirdropAccounts = 5
)

// GenParams returns a randomized uint64 in the range of [1, 10000]
func GenExpiryBlock(r *rand.Rand) uint64 {
	return uint64(1 + r.Intn(10000))
}

// GenDelegationRequirement returns a randomized math.LegacyDec in the range of [0.0, 1.0]
func GenDelegationRequirement(r *rand.Rand) math.LegacyDec {
	return math.LegacyNewDecWithPrec(int64(r.Intn(100)), 1)
}

// GenAirdropFactor returns a randomized math.LegacyDec in the range of [0.0, 1.0]
func GenAirdropFactor(r *rand.Rand) math.LegacyDec {
	return math.LegacyNewDecWithPrec(int64(r.Intn(100)), 1)
}

func RandomizedGenState(simState *module.SimulationState) {
	airdrdopGenesis := types.DefaultGenesisState()

	var expiryBlock uint64
	simState.AppParams.GetOrGenerate(
		expiryBlockKey, &expiryBlock, simState.Rand,
		func(r *rand.Rand) { expiryBlock = GenExpiryBlock(r) },
	)

	var delegationRequirement math.LegacyDec
	simState.AppParams.GetOrGenerate(
		delegationRequirementKey, &delegationRequirement, simState.Rand,
		func(r *rand.Rand) { delegationRequirement = GenDelegationRequirement(r) },
	)

	var airdropFactor math.LegacyDec
	simState.AppParams.GetOrGenerate(
		airdropFactorKey, &airdropFactor, simState.Rand,
		func(r *rand.Rand) { airdropFactor = GenAirdropFactor(r) },
	)

	airdrdopGenesis.Params = types.Params{
		ExpiryBlock:           expiryBlock,
		DelegationRequirement: &delegationRequirement,
		AirdropFactor:         &airdropFactor,
	}

	bz, err := json.MarshalIndent(&airdrdopGenesis.Params, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated airdrop parameters:\n%s\n", bz)

	airdrdopGenesis.AirdropAccounts = GenerateAirdropAccounts(simState)

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(airdrdopGenesis)
}

// Create new account addressees for airdrop accounts because we need to create new
// vesting accounts and the bank module initializes all the existing addresses
func GenerateAirdropAccounts(simState *module.SimulationState) []*types.AirdropAccount {
	accs := simulation.RandomAccounts(simState.Rand, numAirdropAccounts)

	startTime := simState.GenTimestamp.Unix()
	airdropAccounts := make([]*types.AirdropAccount, numAirdropAccounts)
	for i, acc := range accs {
		vestingEndTime := int64(simulation.RandIntBetween(simState.Rand, int(startTime)+1, int(startTime+(60*60*12))))
		airdropAccounts[i] = types.NewAirdropAccount(
			acc.Address.String(),
			uint64(10000),
			vestingEndTime,
		)
	}
	return airdropAccounts
}
