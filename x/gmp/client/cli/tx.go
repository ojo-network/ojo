package cli

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
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
		GetCmdRelayWithContractCall(),
		GetCmdCreatePayment(),
	)

	return cmd
}

func GetCmdRelay() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `relay [destination-chain] [ojo-contract-address] [denoms] [amount]`,
		Args:  cobra.ExactArgs(5),
		Short: "Relay oracle data via Axelar GMP",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if args[0] == "" {
				return fmt.Errorf("destination-chain cannot be empty")
			}
			if args[1] == "" {
				return fmt.Errorf("ojo-contract-address cannot be empty")
			}
			if args[2] == "" {
				return fmt.Errorf("denoms cannot be empty")
			}

			tokens := sdk.Coin{}
			// normalize the coin denom
			if args[4] != "" {
				coin, err := sdk.ParseCoinNormalized(args[4])
				if err != nil {
					return err
				}
				if !strings.HasPrefix(coin.Denom, "ibc/") {
					denomTrace := ibctransfertypes.ParseDenomTrace(coin.Denom)
					coin.Denom = denomTrace.IBCDenom()
				}
				tokens = coin
			}

			// convert denoms to string array
			denoms := strings.Split(args[2], ",")

			commandSelector, err := base64.StdEncoding.DecodeString("")
			if err != nil {
				return err
			}
			commandParams, err := base64.StdEncoding.DecodeString("")
			if err != nil {
				return err
			}

			msg := types.NewMsgRelay(
				clientCtx.GetFromAddress().String(),
				args[0],                // destination-chain e.g. "Ethereum"
				args[1],                // ojo-contract-address e.g. "0x001"
				"",                     // customer-contract-address e.g. "0x002"
				tokens,                 // amount
				denoms,                 // denoms
				commandSelector,        // command-selector
				commandParams,          // command-params
				0,
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

func GetCmdRelayWithContractCall() *cobra.Command {
	cmd := &cobra.Command{
		Use: `relay-with-contract-call [destination-chain] [ojo-contract-address] [client-contract-address] ` +
			`[command-selector] [command-params] [denoms] [amount]`,
		Args:  cobra.ExactArgs(8),
		Short: "Relay oracle data via Axelar GMP and call contract method with oracle data",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if args[0] == "" {
				return fmt.Errorf("destination-chain cannot be empty")
			}
			if args[1] == "" {
				return fmt.Errorf("ojo-contract-address cannot be empty")
			}
			if args[2] == "" {
				return fmt.Errorf("client-contract-address cannot be empty")
			}
			if args[3] == "" {
				return fmt.Errorf("command-selector cannot be empty")
			}
			if args[4] == "" {
				return fmt.Errorf("command-params cannot be empty")
			}
			if args[5] == "" {
				return fmt.Errorf("denoms cannot be empty")
			}

			tokens := sdk.Coin{}
			// normalize the coin denom
			if args[7] != "" {
				coin, err := sdk.ParseCoinNormalized(args[7])
				if err != nil {
					return err
				}
				if !strings.HasPrefix(coin.Denom, "ibc/") {
					denomTrace := ibctransfertypes.ParseDenomTrace(coin.Denom)
					coin.Denom = denomTrace.IBCDenom()
				}
				tokens = coin
			}

			// convert denoms to string array
			denoms := strings.Split(args[5], ",")

			commandSelector, err := base64.StdEncoding.DecodeString(args[3])
			if err != nil {
				return err
			}
			commandParams, err := base64.StdEncoding.DecodeString(args[4])
			if err != nil {
				return err
			}

			msg := types.NewMsgRelay(
				clientCtx.GetFromAddress().String(),
				args[0],                // destination-chain e.g. "Ethereum"
				args[1],                // ojo-contract-address e.g. "0x001"
				args[2],                // customer-contract-address e.g. "0x002"
				tokens,                 // amount
				denoms,                 // denoms
				commandSelector,        // command-selector
				commandParams,          // command-params
				0,
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

func GetCmdCreatePayment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-payment [destination-chain] [denom] [deviation] [heartbeat] [token]",
		Args:  cobra.ExactArgs(5),
		Short: "Create a payment to relay price data",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			deviation, err := math.LegacyNewDecFromStr(args[2])
			if err != nil {
				return err
			}
			heartbeat, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return err
			}

			tokens := sdk.Coin{}
			// normalize the coin denom
			if args[4] != "" {
				coin, err := sdk.ParseCoinNormalized(args[4])
				if err != nil {
					return err
				}
				if !strings.HasPrefix(coin.Denom, "ibc/") {
					denomTrace := ibctransfertypes.ParseDenomTrace(coin.Denom)
					coin.Denom = denomTrace.IBCDenom()
				}
				tokens = coin
			}

			msg := types.NewMsgCreatePayment(
				clientCtx.GetFromAddress().String(),
				args[0],
				args[1],
				tokens,
				deviation,
				heartbeat,
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
