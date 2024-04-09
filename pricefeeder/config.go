package pricefeeder

import (
	"fmt"
	"time"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cast"
)

const (
	DefaultConfigTemplate = `
[pricefeeder]
# Path to price feeder config file
config_path = ""
# Specifies whether the currency pair providers and currency deviation threshold values should
# be read from the oracle module's on chain parameters or the price feeder config
chain_config = false
# Log level of price feeder process
log_level = "info"
# Time interval that the price feeder's oracle process waits before fetching for new prices
oracle_tick_time = "5s"
`
)

const (
	FlagConfigPath     = "pricefeeder.config_path"
	FlagChainConfig    = "pricefeeder.chain_config"
	FlagLogLevel       = "pricefeeder.log_level"
	FlagOracleTickTime = "pricefeeder.oracle_tick_time"
)

// AppConfig defines the app configuration for the price feeder that must be set in the app.toml file.
type AppConfig struct {
	ConfigPath     string        `mapstructure:"config_path"`
	ChainConfig    bool          `mapstructure:"chain_config"`
	LogLevel       string        `mapstructure:"log_level"`
	OracleTickTime time.Duration `mapstructure:"oracle_tick_time"`
}

// ValidateBasic performs basic validation of the price feeder app config.
func (c *AppConfig) ValidateBasic() error {
	if c.ConfigPath == "" {
		return fmt.Errorf("path to price feeder config must be set")
	}

	if c.OracleTickTime <= 0 {
		return fmt.Errorf("oracle tick time must be greater than 0")
	}

	return nil
}

// ReadConfigFromAppOpts reads the config parameters from the AppOptions and returns the config.
func ReadConfigFromAppOpts(opts servertypes.AppOptions) (AppConfig, error) {
	var (
		cfg AppConfig
		err error
	)

	if v := opts.Get(FlagConfigPath); v != nil {
		if cfg.ConfigPath, err = cast.ToStringE(v); err != nil {
			return cfg, err
		}
	}

	if v := opts.Get(FlagChainConfig); v != nil {
		if cfg.ChainConfig, err = cast.ToBoolE(v); err != nil {
			return cfg, err
		}
	}

	if v := opts.Get(FlagLogLevel); v != nil {
		if cfg.LogLevel, err = cast.ToStringE(v); err != nil {
			return cfg, err
		}
	}

	if v := opts.Get(FlagOracleTickTime); v != nil {
		if cfg.OracleTickTime, err = cast.ToDurationE(v); err != nil {
			return cfg, err
		}
	}

	if err := cfg.ValidateBasic(); err != nil {
		return cfg, err
	}

	return cfg, err
}
