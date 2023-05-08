package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (aa *AirdropAccount) OriginAccAddress() sdk.AccAddress {
	return sdk.AccAddress(aa.OriginAddress)
}

func (aa *AirdropAccount) ClaimAccAddress() sdk.AccAddress {
	return sdk.AccAddress(aa.ClaimAddress)
}

func (aa AirdropAccount) OriginCoins() sdk.Coins {
	return sdk.NewCoins(sdk.NewCoin("ojo", sdk.NewIntFromUint64(aa.OriginAmount)))
}

func (aa *AirdropAccount) ClaimCoins() sdk.Coins {
	return sdk.NewCoins(sdk.NewCoin("ojo", sdk.NewIntFromUint64(aa.ClaimAmount)))
}

func (aa *AirdropAccount) VerifyNotClaimed() error {
	if aa.ClaimAddress != "" {
		return errors.Wrapf(
			ErrAirdropAlreadyClaimed,
			"already claimed by address %s",
			aa.ClaimAddress,
		)
	}
	return nil
}
