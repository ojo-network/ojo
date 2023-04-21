package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	DefaultExpiryBlock           = uint64(5000)
	DefaultDelegationRequirement = sdk.NewDecWithPrec(1, 1)
	DefaultAirdropFactor         = sdk.NewDecWithPrec(1, 1)
)

func DefaultParams() Params {
	return Params{
		ExpiryBlock:           DefaultExpiryBlock,
		DelegationRequirement: DefaultDelegationRequirement,
		AirdropFactor:         DefaultAirdropFactor,
	}
}
