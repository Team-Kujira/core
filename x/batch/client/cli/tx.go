package cli

import (
	"fmt"
	"strings"

	"cosmossdk.io/math"
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
		NewBatchResetDelegationCmd(),
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

// NewBatchResetDelegationCmd returns a CLI command handler for creating a MsgBatchResetDelegation transaction.
func NewBatchResetDelegationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-reset-delegation",
		Short: "Reset delegations in batch for a delegator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Reset delegations in batch for a single delegator.

Example:
$ %[1]s tx batch batch-reset-delegation kujiravaloper1uf9knclap4a0vrqn8dj726anyhst7l0253g0da,kujiravaloper1ujhlm5qxyt2hn5fxq8wll805tsxcamqfhsty9a 100,200 --from mykey
`,
				version.AppName, flags.FlagBroadcastMode, flags.BroadcastSync, flags.BroadcastAsync,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			delAddr := clientCtx.GetFromAddress()

			valAddrs := strings.Split(args[0], ",")
			amountStrs := strings.Split(args[1], ",")
			amounts := []math.Int{}
			for _, amountStr := range amountStrs {
				amount, ok := math.NewIntFromString(amountStr)
				if !ok {
					return types.ErrInvalidAmount
				}
				amounts = append(amounts, amount)
			}
			// The transaction cannot be generated offline since it requires a query
			// to get all the validators.
			if clientCtx.Offline {
				return fmt.Errorf("cannot generate tx in offline mode")
			}

			msgs := []sdk.Msg{types.NewMsgBatchResetDelegation(delAddr, valAddrs, amounts)}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msgs...)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
