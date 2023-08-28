package app

import (
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ojo-network/ojo/app/params"
)

// MakeEncodingConfig returns the application's encoding configuration with all
// types and interfaces registered.
func MakeEncodingConfig() testutil.TestEncodingConfig {
	encodingConfig := params.MakeEncodingConfig()
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	// register auth type interfaces
	authtypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	return encodingConfig
}
