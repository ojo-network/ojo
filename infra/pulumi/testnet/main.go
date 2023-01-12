package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{
		DisableHTMLEscape: true,
	})
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")
		var network Network
		var secrets []NodeSecretConfig
		cfg.RequireObject("config", &network)

		return network.Provision(ctx, secrets)
	})
}
