package cli

import (
	"fmt"
	"strings"

	"github.com/Team-Kujira/core/x/oracle/types"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	oracleTxCmd := &cobra.Command{
		Use:                        "oracle",
		Short:                      "Oracle transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	oracleTxCmd.AddCommand(
		GetCmdDelegateFeederPermission(),
		GetCmdAggregateExchangeRatePrevote(),
		GetCmdAggregateExchangeRateVote(),
	)

	return oracleTxCmd
}

// GetCmdDelegateFeederPermission will create a feeder permission delegation tx and sign it with the given key.
func GetCmdDelegateFeederPermission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-feeder [feeder]",
		Args:  cobra.ExactArgs(1),
		Short: "Delegate the permission to vote for the oracle to an address",
		Long: strings.TrimSpace(`
Delegate the permission to submit exchange rate votes for the oracle to an address.

Delegation can keep your validator operator key offline and use a separate replaceable key online.

$ kujirad tx oracle set-feeder kujira1...

where "kujira1..." is the address you want to delegate your voting rights to.
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Get from address
			voter := clientCtx.GetFromAddress()

			// The address the right is being delegated from
			validator := sdk.ValAddress(voter)

			feederStr := args[0]
			feeder, err := sdk.AccAddressFromBech32(feederStr)
			if err != nil {
				return err
			}

			msgs := []sdk.Msg{types.NewMsgDelegateFeedConsent(validator, feeder)}
			for _, msg := range msgs {
				if err := msg.ValidateBasic(); err != nil {
					return err
				}
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msgs...)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdAggregateExchangeRatePrevote will create a aggregateExchangeRatePrevote tx and sign it with the given key.
func GetCmdAggregateExchangeRatePrevote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aggregate-prevote [salt] [exchange-rates] [validator]",
		Args:  cobra.RangeArgs(2, 3),
		Short: "Submit an oracle aggregate prevote for the exchange rates",
		Long: strings.TrimSpace(`
Submit an oracle aggregate prevote for the exchange rates of multiple denoms.
The purpose of aggregate prevote is to hide aggregate exchange rate vote with hash which is formatted 
as hex string in SHA256("{salt}:{exchange_rate}{denom},...,{exchange_rate}{denom}:{voter}")

# Aggregate Prevote
$ kujirad tx oracle aggregate-prevote 1234 0.1ATOM,1.001USDT

where "ATOM,USDT" is the denominating currencies, and "0.1.0,1.001" is the exchange rates of USD from the voter's point of view.

If voting from a voting delegate, set "validator" to the address of the validator to vote on behalf of:
$ kujirad tx oracle aggregate-prevote 1234 0.1ATOM,1.001USDT kujiravaloper1...
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			salt := args[0]
			exchangeRatesStr := args[1]
			_, err = types.ParseExchangeRateTuples(exchangeRatesStr)
			if err != nil {
				return fmt.Errorf("given exchange_rates {%s} is not a valid format; exchange_rate should be formatted as DecCoins; %s", exchangeRatesStr, err.Error())
			}

			// Get from address
			voter := clientCtx.GetFromAddress()

			// By default the voter is voting on behalf of itself
			validator := sdk.ValAddress(voter)

			// Override validator if validator is given
			if len(args) == 3 {
				parsedVal, err := sdk.ValAddressFromBech32(args[2])
				if err != nil {
					return errorsmod.Wrap(err, "validator address is invalid")
				}
				validator = parsedVal
			}

			hash := types.GetAggregateVoteHash(salt, exchangeRatesStr, validator)
			msgs := []sdk.Msg{types.NewMsgAggregateExchangeRatePrevote(hash, voter, validator)}
			for _, msg := range msgs {
				if err := msg.ValidateBasic(); err != nil {
					return err
				}
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msgs...)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdAggregateExchangeRateVote will create a aggregateExchangeRateVote tx and sign it with the given key.
func GetCmdAggregateExchangeRateVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aggregate-vote [salt] [exchange-rates] [validator]",
		Args:  cobra.RangeArgs(2, 3),
		Short: "Submit an oracle aggregate vote for the exchange_rates of Luna",
		Long: strings.TrimSpace(`
Submit a aggregate vote for the exchange_rates of Luna w.r.t the input denom. Companion to a prevote submitted in the previous vote period. 

$ kujirad tx oracle aggregate-vote 1234 0.1ATOM,1.001USDT 

where "ATOM,USDT" is the denominating currencies, and "0.1.0,1.001" is the exchange rates of USD from the voter's point of view.

"salt" should match the salt used to generate the SHA256 hex in the aggregated pre-vote. 

If voting from a voting delegate, set "validator" to the address of the validator to vote on behalf of:
$ kujirad tx oracle aggregate-vote 1234 0.1ATOM,1.001USDT kujiravaloper1....
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			salt := args[0]
			exchangeRatesStr := args[1]
			_, err = types.ParseExchangeRateTuples(exchangeRatesStr)
			if err != nil {
				return fmt.Errorf("given exchange_rate {%s} is not a valid format; exchange rate should be formatted as DecCoin; %s", exchangeRatesStr, err.Error())
			}

			// Get from address
			voter := clientCtx.GetFromAddress()

			// By default the voter is voting on behalf of itself
			validator := sdk.ValAddress(voter)

			// Override validator if validator is given
			if len(args) == 3 {
				parsedVal, err := sdk.ValAddressFromBech32(args[2])
				if err != nil {
					return errorsmod.Wrap(err, "validator address is invalid")
				}
				validator = parsedVal
			}

			msgs := []sdk.Msg{types.NewMsgAggregateExchangeRateVote(salt, exchangeRatesStr, voter, validator)}
			for _, msg := range msgs {
				if err := msg.ValidateBasic(); err != nil {
					return err
				}
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msgs...)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
