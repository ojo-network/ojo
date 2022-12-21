package main

import (
	"os"
	"strings"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	ojoapp "github.com/ojo-network/ojo/app"
	"github.com/ojo-network/ojo/app/params"
	"github.com/ojo-network/ojo/cmd/ojod/cmd"
)

func main() {
	params.SetAddressPrefixes()
	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, strings.ToUpper(params.Name), ojoapp.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
