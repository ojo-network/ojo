package params

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
)

// MakeEncodingConfig creates an EncodingConfig for Amino-based tests.
func MakeEncodingConfig(modules ...module.AppModuleBasic) testutil.TestEncodingConfig {
	encCfg := testutil.MakeTestEncodingConfig(modules...)

	// register auth type interfaces
	authtypes.RegisterInterfaces(encCfg.InterfaceRegistry)

	return encCfg
}
