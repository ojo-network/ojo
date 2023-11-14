package cli

import (
	"fmt"
	"strconv"
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
		Use: `relay [destination-chain] [destination-address] [command-selector] ` +
			`[command-params] [timestamp] [denoms] [amount]`,
		Args:  cobra.ExactArgs(4),
		Short: "Relay oracle data via Axelar GMP",
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
			if args[2] == "" {
				return fmt.Errorf("command-selector cannot be empty")
			}
			if args[3] == "" {
				return fmt.Errorf("command-params cannot be empty")
			}
			if args[4] == "" {
				return fmt.Errorf("timestamp cannot be empty")
			}
			if args[5] == "" {
				return fmt.Errorf("denoms cannot be empty")
			}
			if args[6] == "" {
				return fmt.Errorf("amount cannot be empty")
			}

			// normalize the coin denom
			coin, err := sdk.ParseCoinNormalized(args[6])
			if err != nil {
				return err
			}
			if !strings.HasPrefix(coin.Denom, "ibc/") {
				denomTrace := ibctransfertypes.ParseDenomTrace(coin.Denom)
				coin.Denom = denomTrace.IBCDenom()
			}

			// convert denoms to string array
			denoms := strings.Split(args[3], ",")

			// convert timestamp string to int64
			timestamp, err := strconv.ParseInt(args[4], 10, 64)
			if err != nil {
				return err
			}

			// convert command-selector to []byte
			var commandSelector []byte
			copy(commandSelector, args[2])

			// convert command-params to []byte
			var commandParams []byte
			copy(commandParams, args[3])

			msg := types.NewMsgRelay(
				clientCtx.GetFromAddress().String(),
				args[0],         // destination-chain
				args[1],         // destination-address
				coin,            // amount
				denoms,          // denoms
				commandSelector, // command-selector
				commandParams,   // command-params
				timestamp,       // timestamp
			)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
