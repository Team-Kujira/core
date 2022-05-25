package cli

import (
	"strconv"

	"kujira/x/scheduler/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdCreateHook() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-hook [contract] [executor] [msg] [frequency]",
		Short: "Create a new hook",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argContract := args[0]
			argExecutor := args[1]
			argMsg := args[2]
			argFrequency, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateHook(clientCtx.GetFromAddress().String(), argContract, argExecutor, []byte(argMsg), argFrequency)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdateHook() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-hook [id] [contract] [executor] [msg] [frequency]",
		Short: "Update a hook",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			argContract := args[1]

			argExecutor := args[2]

			argMsg := args[3]

			argFrequency, err := strconv.ParseInt(args[4], 10, 64)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateHook(clientCtx.GetFromAddress().String(), id, argContract, argExecutor, []byte(argMsg), argFrequency)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdDeleteHook() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-hook [id]",
		Short: "Delete a hook by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteHook(clientCtx.GetFromAddress().String(), id)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
