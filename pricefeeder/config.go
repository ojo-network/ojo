package pricefeeder

import (
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cast"
)

const (
	DefaultConfigTemplate = `
[pricefeeder]
# Log level of price feeder process.
log_level = "info"

# Enable the price feeder.
enable = true
`
)

const (
	FlagLogLevel          = "pricefeeder.log_level"
	FlagEnablePriceFeeder = "pricefeeder.enable"
)

// AppConfig defines the app configuration for the price feeder that must be set in the app.toml file.
type AppConfig struct {
	LogLevel string `mapstructure:"log_level"`
	Enable   bool   `mapstructure:"enable"`
}

// ReadConfigFromAppOpts reads the config parameters from the AppOptions and returns the config.
func ReadConfigFromAppOpts(opts servertypes.AppOptions) (AppConfig, error) {
	var (
		cfg AppConfig
		err error
	)

	cfg.LogLevel = "info"
	cfg.Enable = true

	// Override with values from AppOptions if provided
	if v := opts.Get(FlagLogLevel); v != nil {
		if cfg.LogLevel, err = cast.ToStringE(v); err != nil {
			return cfg, err
		}
	}

	if v := opts.Get(FlagEnablePriceFeeder); v != nil {
		if cfg.Enable, err = cast.ToBoolE(v); err != nil {
			return cfg, err
		}
	}

	return cfg, err
}
