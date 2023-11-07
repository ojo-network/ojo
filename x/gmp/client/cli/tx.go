package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	"github.com/ojo-network/ojo/x/gmp/types"
	"github.com/spf13/cobra"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Transaction commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetCmdRelay(),
	)

	return cmd
}

func GetCmdRelay() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relay [destination-chain] [destination-address] [amount] [comma-separated list of tokens]",
		Args:  cobra.ExactArgs(4),
		Short: "Relay via axelar GMP to an address",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Flags().Set(flags.FlagFrom, args[0]); err != nil {
				return err
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if args[0] == "" {
				return fmt.Errorf("destination-chain cannot be empty")
			}
			if args[1] == "" {
				return fmt.Errorf("destination-address cannot be empty")
			}
			if args[3] == "" {
				return fmt.Errorf("denoms cannot be empty")
			}

			// Normalize the coin denom
			coin, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}
			if !strings.HasPrefix(coin.Denom, "ibc/") {
				denomTrace := ibctransfertypes.ParseDenomTrace(coin.Denom)
				coin.Denom = denomTrace.IBCDenom()
			}

			denoms := strings.Split(args[3], ",")

			msg := types.NewMsgRelay(clientCtx.GetFromAddress().String(), args[0], args[1], coin, denoms)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
