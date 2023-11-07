package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ojo-network/ojo/client/tx"
)

// CreateAccount creates a new account with a random mnemonic and returns the address
func CreateAccount(s *IntegrationTestSuite) sdk.AccAddress {
	mnemonic, err := tx.CreateMnemonic()
	s.Require().NoError(err)
	account, _, err := tx.CreateAccountFromMnemonic("test", mnemonic)
	s.Require().NoError(err)
	address, err := account.GetAddress()
	s.Require().NoError(err)
	return address
}
