package types

func NewGenesisState(
	params Params,
	airdropAccounts []AirdropAccount,
) *GenesisState {
	return &GenesisState{
		Params:          params,
		AirdropAccounts: airdropAccounts,
	}
}

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:          DefaultParams(),
		AirdropAccounts: []AirdropAccount{},
	}
}
