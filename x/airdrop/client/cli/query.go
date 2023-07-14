package cli

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/ojo-network/ojo/util/cli"
	"github.com/ojo-network/ojo/x/airdrop/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQueryAirdropAccount(),
	)

	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current Airdrop params",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(cmd.Context(), &types.ParamsRequest{})
			return cli.PrintOrErr(res, err, clientCtx)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryAirdropAccount implements the query airdrop account command.
func GetCmdQueryAirdropAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "airdrop-account [address]",
		Args:  cobra.ExactArgs(1),
		Short: "Query a specific airdrop account",
		Long: strings.TrimSpace(`
Query a specific airdrop account.

$ ojod query airdrop airdrop-account [address]
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			if _, err = sdk.ValAddressFromBech32(args[0]); err != nil {
				return err
			}

			res, err := queryClient.AirdropAccount(
				cmd.Context(),
				&types.AirdropAccountRequest{
					Address: args[0],
				},
			)
			return cli.PrintOrErr(res, err, clientCtx)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
