package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"

	srvconfig "github.com/cosmos/cosmos-sdk/server/config"
	tmconfig "github.com/tendermint/tendermint/config"
)

func (o *Orchestrator) initOjoConfigs() (dir string, err error) {
}

func initOjoConfig() (dir string, err error) {
	dir, err = os.MkdirTemp("", "e2e-configs")
	if err != nil {
		return
	}

	_, err = copyFile(
		filepath.Join("./config/", "ojo_bootstrap.sh"),
		filepath.Join(dir, "ojo_bootstrap.sh"),
	)
	if err != nil {
		return
	}

	configPath := filepath.Join(dir, "config.toml")
	config := tmconfig.DefaultConfig()
	config.P2P.ListenAddress = "tcp://0.0.0.0:26656"
	config.RPC.ListenAddress = fmt.Sprintf("tcp://0.0.0.0:%s", ojo_tmrpc_port)
	config.StateSync.Enable = false
	config.P2P.AddrBookStrict = false
	config.P2P.Seeds = ""
	tmconfig.WriteConfigFile(configPath, config)

	appCfgPath := filepath.Join(dir, "app.toml")
	appConfig := srvconfig.DefaultConfig()
	appConfig.API.Enable = true
	appConfig.MinGasPrices = ojoMinGasPrice
	srvconfig.WriteConfigFile(appCfgPath, appConfig)

	return
}
