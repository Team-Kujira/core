package cli

import (
	"fmt"
	"strings"

	"github.com/Team-Kujira/core/x/batch/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	distTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Distribution extension commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	distTxCmd.AddCommand(
		NewWithdrawAllRewardsCmd(),
	)

	return distTxCmd
}

// NewWithdrawAllRewardsCmd returns a CLI command handler for creating a MsgWithdrawDelegatorReward transaction.
func NewWithdrawAllRewardsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-all-rewards",
		Short: "withdraw all delegations rewards for a delegator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Withdraw all rewards for a single delegator.

Example:
$ %[1]s tx distribution withdraw-all-rewards --from mykey
`,
				version.AppName, flags.FlagBroadcastMode, flags.BroadcastSync, flags.BroadcastAsync,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			delAddr := clientCtx.GetFromAddress()

			// The transaction cannot be generated offline since it requires a query
			// to get all the validators.
			if clientCtx.Offline {
				return fmt.Errorf("cannot generate tx in offline mode")
			}

			msgs := []sdk.Msg{types.NewMsgWithdrawAllDelegatorRewards(delAddr)}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msgs...)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
