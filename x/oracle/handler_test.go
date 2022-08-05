package oracle_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"kujira/x/oracle/keeper"
	"kujira/x/oracle/types"
)

func TestOracleFilters(t *testing.T) {
	input, h := setup(t)

	// Case 1: non-oracle message being sent fails
	bankMsg := banktypes.MsgSend{}
	_, err := h(input.Ctx, &bankMsg)
	require.Error(t, err)

	// Case 2: Normal MsgAggregateExchangeRatePrevote submission goes through
	salt := "fc5bb0bc63e54b2918d9334bf3259f5dc575e8d7a4df4e836dd80f1ad62aa89b"

	hash := types.GetAggregateVoteHash(salt, randomExchangeRate.String()+types.TestDenomD, keeper.ValAddrs[0])
	prevoteMsg := types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[0], keeper.ValAddrs[0])
	_, err = h(input.Ctx, prevoteMsg)
	require.NoError(t, err)

	// // Case 3: Normal MsgAggregateExchangeRateVote submission goes through keeper.Addrs
	voteMsg := types.NewMsgAggregateExchangeRateVote(salt, randomExchangeRate.String()+types.TestDenomD, keeper.Addrs[0], keeper.ValAddrs[0])
	_, err = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.NoError(t, err)

	// Case 4: a non-validator sending an oracle message fails
	nonValidatorPub := secp256k1.GenPrivKey().PubKey()
	nonValidatorAddr := nonValidatorPub.Address()
	salt = "4c81c928f466a08b07171def7aeb2b3c266df7bb7486158a15a2291a7d55c8f9"
	hash = types.GetAggregateVoteHash(salt, randomExchangeRate.String()+types.TestDenomD, sdk.ValAddress(nonValidatorAddr))

	prevoteMsg = types.NewMsgAggregateExchangeRatePrevote(hash, sdk.AccAddress(nonValidatorAddr), sdk.ValAddress(nonValidatorAddr))
	_, err = h(input.Ctx, prevoteMsg)
	require.Error(t, err)
}

func TestFeederDelegation(t *testing.T) {
	input, h := setup(t)

	salt := "fc5bb0bc63e54b2918d9334bf3259f5dc575e8d7a4df4e836dd80f1ad62aa89b"
	hash := types.GetAggregateVoteHash(salt, randomExchangeRate.String()+types.TestDenomD, keeper.ValAddrs[0])

	// Case 1: empty message
	delegateFeedConsentMsg := types.MsgDelegateFeedConsent{}
	_, err := h(input.Ctx, &delegateFeedConsentMsg)
	require.Error(t, err)

	// Case 2: Normal Prevote - without delegation
	prevoteMsg := types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[0], keeper.ValAddrs[0])
	_, err = h(input.Ctx, prevoteMsg)
	require.NoError(t, err)

	// Case 2.1: Normal Prevote - with delegation fails
	prevoteMsg = types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[1], keeper.ValAddrs[0])
	_, err = h(input.Ctx, prevoteMsg)
	require.Error(t, err)

	// Case 2.2: Normal Vote - without delegation
	voteMsg := types.NewMsgAggregateExchangeRateVote(salt, randomExchangeRate.String()+types.TestDenomD, keeper.Addrs[0], keeper.ValAddrs[0])
	_, err = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.NoError(t, err)

	// Case 2.3: Normal Vote - with delegation fails
	voteMsg = types.NewMsgAggregateExchangeRateVote(salt, randomExchangeRate.String()+types.TestDenomD, keeper.Addrs[1], keeper.ValAddrs[0])
	_, err = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.Error(t, err)

	// Case 3: Normal MsgDelegateFeedConsent succeeds
	msg := types.NewMsgDelegateFeedConsent(keeper.ValAddrs[0], keeper.Addrs[1])
	_, err = h(input.Ctx, msg)
	require.NoError(t, err)

	// Case 4.1: Normal Prevote - without delegation fails
	prevoteMsg = types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[2], keeper.ValAddrs[0])
	_, err = h(input.Ctx, prevoteMsg)
	require.Error(t, err)

	// Case 4.2: Normal Prevote - with delegation succeeds
	prevoteMsg = types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[1], keeper.ValAddrs[0])
	_, err = h(input.Ctx, prevoteMsg)
	require.NoError(t, err)

	// Case 4.3: Normal Vote - without delegation fails
	voteMsg = types.NewMsgAggregateExchangeRateVote(salt, randomExchangeRate.String()+types.TestDenomD, keeper.Addrs[2], keeper.ValAddrs[0])
	_, err = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.Error(t, err)

	// Case 4.4: Normal Vote - with delegation succeeds
	voteMsg = types.NewMsgAggregateExchangeRateVote(salt, randomExchangeRate.String()+types.TestDenomD, keeper.Addrs[1], keeper.ValAddrs[0])
	_, err = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.NoError(t, err)
}

func TestAggregatePrevoteVote(t *testing.T) {
	input, h := setup(t)

	salt := "fc5bb0bc63e54b2918d9334bf3259f5dc575e8d7a4df4e836dd80f1ad62aa89b"
	exchangeRatesStr := fmt.Sprintf("1000.23%s,0.29%s,0.27%s", types.TestDenomC, types.TestDenomB, types.TestDenomD)
	otherExchangeRateStr := fmt.Sprintf("1000.12%s,0.29%s,0.27%s", types.TestDenomC, types.TestDenomB, types.TestDenomB)
	unintendedExchageRateStr := fmt.Sprintf("1000.23%s,0.29%s,0.27%s", types.TestDenomC, types.TestDenomB, types.TestDenomE)
	invalidExchangeRateStr := fmt.Sprintf("1000.23%s,0.29%s,0.27", types.TestDenomC, types.TestDenomB)

	hash := types.GetAggregateVoteHash(salt, exchangeRatesStr, keeper.ValAddrs[0])

	aggregateExchangeRatePrevoteMsg := types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[0], keeper.ValAddrs[0])
	_, err := h(input.Ctx, aggregateExchangeRatePrevoteMsg)
	require.NoError(t, err)

	// Unauthorized feeder
	aggregateExchangeRatePrevoteMsg = types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[1], keeper.ValAddrs[0])
	_, err = h(input.Ctx, aggregateExchangeRatePrevoteMsg)
	require.Error(t, err)

	// Invalid reveal period
	aggregateExchangeRateVoteMsg := types.NewMsgAggregateExchangeRateVote(salt, exchangeRatesStr, sdk.AccAddress(keeper.Addrs[0]), keeper.ValAddrs[0])
	_, err = h(input.Ctx, aggregateExchangeRateVoteMsg)
	require.Error(t, err)

	// Invalid reveal period
	input.Ctx = input.Ctx.WithBlockHeight(2)
	aggregateExchangeRateVoteMsg = types.NewMsgAggregateExchangeRateVote(salt, exchangeRatesStr, sdk.AccAddress(keeper.Addrs[0]), keeper.ValAddrs[0])
	_, err = h(input.Ctx, aggregateExchangeRateVoteMsg)
	require.Error(t, err)

	// Other exchange rate with valid real period
	input.Ctx = input.Ctx.WithBlockHeight(1)
	aggregateExchangeRateVoteMsg = types.NewMsgAggregateExchangeRateVote(salt, otherExchangeRateStr, sdk.AccAddress(keeper.Addrs[0]), keeper.ValAddrs[0])
	_, err = h(input.Ctx, aggregateExchangeRateVoteMsg)
	require.Error(t, err)

	// Invalid exchange rate with valid real period
	input.Ctx = input.Ctx.WithBlockHeight(1)
	aggregateExchangeRateVoteMsg = types.NewMsgAggregateExchangeRateVote(salt, invalidExchangeRateStr, sdk.AccAddress(keeper.Addrs[0]), keeper.ValAddrs[0])
	_, err = h(input.Ctx, aggregateExchangeRateVoteMsg)
	require.Error(t, err)

	// Unauthorized feeder
	aggregateExchangeRateVoteMsg = types.NewMsgAggregateExchangeRateVote(salt, invalidExchangeRateStr, sdk.AccAddress(keeper.Addrs[1]), keeper.ValAddrs[0])
	_, err = h(input.Ctx, aggregateExchangeRateVoteMsg)
	require.Error(t, err)

	// Unintended denom vote
	aggregateExchangeRateVoteMsg = types.NewMsgAggregateExchangeRateVote(salt, unintendedExchageRateStr, sdk.AccAddress(keeper.Addrs[0]), keeper.ValAddrs[0])
	_, err = h(input.Ctx, aggregateExchangeRateVoteMsg)
	require.Error(t, err)

	// Valid exchange rate reveal submission
	input.Ctx = input.Ctx.WithBlockHeight(1)
	aggregateExchangeRateVoteMsg = types.NewMsgAggregateExchangeRateVote(salt, exchangeRatesStr, sdk.AccAddress(keeper.Addrs[0]), keeper.ValAddrs[0])
	_, err = h(input.Ctx, aggregateExchangeRateVoteMsg)
	require.NoError(t, err)
}
