package oracle_test

import (
	"fmt"
	"math"
	"sort"
	"testing"

	"github.com/cometbft/cometbft/libs/rand"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/Team-Kujira/core/x/oracle"
	"github.com/Team-Kujira/core/x/oracle/keeper"
	"github.com/Team-Kujira/core/x/oracle/types"
)

func TestOracleThreshold(t *testing.T) {
	input, h := setup(t)
	exchangeRateStr := randomExchangeRate.String() + types.TestDenomD

	// Case 1.
	// Less than the threshold signs, exchange rate consensus fails
	salt := "fc5bb0bc63e54b2918d9334bf3259f5dc575e8d7a4df4e836dd80f1ad62aa89b"
	hash := types.GetAggregateVoteHash(salt, exchangeRateStr, keeper.ValAddrs[0])
	prevoteMsg := types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[0], keeper.ValAddrs[0])
	voteMsg := types.NewMsgAggregateExchangeRateVote(salt, exchangeRateStr, keeper.Addrs[0], keeper.ValAddrs[0])

	_, err1 := h.AggregateExchangeRatePrevote(input.Ctx.WithBlockHeight(0), prevoteMsg)
	_, err2 := h.AggregateExchangeRateVote(input.Ctx.WithBlockHeight(1), voteMsg)
	require.NoError(t, err1)
	require.NoError(t, err2)

	oracle.EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	_, err := input.OracleKeeper.GetExchangeRate(input.Ctx.WithBlockHeight(1), types.TestDenomD)
	require.Error(t, err)

	// Case 2.
	// More than the threshold signs, exchange rate consensus succeeds
	salt = "fc5bb0bc63e54b2918d9334bf3259f5dc575e8d7a4df4e836dd80f1ad62aa89b"
	hash = types.GetAggregateVoteHash(salt, exchangeRateStr, keeper.ValAddrs[0])
	prevoteMsg = types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[0], keeper.ValAddrs[0])
	voteMsg = types.NewMsgAggregateExchangeRateVote(salt, exchangeRateStr, keeper.Addrs[0], keeper.ValAddrs[0])

	_, err1 = h.AggregateExchangeRatePrevote(input.Ctx.WithBlockHeight(0), prevoteMsg)
	_, err2 = h.AggregateExchangeRateVote(input.Ctx.WithBlockHeight(1), voteMsg)
	require.NoError(t, err1)
	require.NoError(t, err2)

	salt = "4c81c928f466a08b07171def7aeb2b3c266df7bb7486158a15a2291a7d55c8f9"
	hash = types.GetAggregateVoteHash(salt, exchangeRateStr, keeper.ValAddrs[1])
	prevoteMsg = types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[1], keeper.ValAddrs[1])
	voteMsg = types.NewMsgAggregateExchangeRateVote(salt, exchangeRateStr, keeper.Addrs[1], keeper.ValAddrs[1])

	_, err1 = h.AggregateExchangeRatePrevote(input.Ctx.WithBlockHeight(0), prevoteMsg)
	_, err2 = h.AggregateExchangeRateVote(input.Ctx.WithBlockHeight(1), voteMsg)
	require.NoError(t, err1)
	require.NoError(t, err2)

	salt = "fc246cf5a18c7a650a6a226ebc589d49a9a814d6f1f586405e8726e5cf2a7d80"
	hash = types.GetAggregateVoteHash(salt, exchangeRateStr, keeper.ValAddrs[2])
	prevoteMsg = types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[2], keeper.ValAddrs[2])
	voteMsg = types.NewMsgAggregateExchangeRateVote(salt, exchangeRateStr, keeper.Addrs[2], keeper.ValAddrs[2])

	_, err1 = h.AggregateExchangeRatePrevote(input.Ctx.WithBlockHeight(0), prevoteMsg)
	_, err2 = h.AggregateExchangeRateVote(input.Ctx.WithBlockHeight(1), voteMsg)
	require.NoError(t, err1)
	require.NoError(t, err2)

	oracle.EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	rate, err := input.OracleKeeper.GetExchangeRate(input.Ctx.WithBlockHeight(1), types.TestDenomD)
	require.NoError(t, err)
	require.Equal(t, randomExchangeRate, rate)

	// Case 3.
	// Increase voting power of absent validator, exchange rate consensus fails
	val, _ := input.StakingKeeper.GetValidator(input.Ctx, keeper.ValAddrs[2])
	input.StakingKeeper.Delegate(input.Ctx.WithBlockHeight(0), keeper.Addrs[2], stakingAmt.MulRaw(3), stakingtypes.Unbonded, val, false)

	salt = "fc5bb0bc63e54b2918d9334bf3259f5dc575e8d7a4df4e836dd80f1ad62aa89b"
	hash = types.GetAggregateVoteHash(salt, exchangeRateStr, keeper.ValAddrs[0])
	prevoteMsg = types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[0], keeper.ValAddrs[0])
	voteMsg = types.NewMsgAggregateExchangeRateVote(salt, exchangeRateStr, keeper.Addrs[0], keeper.ValAddrs[0])

	_, err1 = h.AggregateExchangeRatePrevote(input.Ctx.WithBlockHeight(0), prevoteMsg)
	_, err2 = h.AggregateExchangeRateVote(input.Ctx.WithBlockHeight(1), voteMsg)
	require.NoError(t, err1)
	require.NoError(t, err2)

	salt = "4c81c928f466a08b07171def7aeb2b3c266df7bb7486158a15a2291a7d55c8f9"
	hash = types.GetAggregateVoteHash(salt, exchangeRateStr, keeper.ValAddrs[1])
	prevoteMsg = types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[1], keeper.ValAddrs[1])
	voteMsg = types.NewMsgAggregateExchangeRateVote(salt, exchangeRateStr, keeper.Addrs[1], keeper.ValAddrs[1])

	_, err1 = h.AggregateExchangeRatePrevote(input.Ctx.WithBlockHeight(0), prevoteMsg)
	_, err2 = h.AggregateExchangeRateVote(input.Ctx.WithBlockHeight(1), voteMsg)
	require.NoError(t, err1)
	require.NoError(t, err2)

	oracle.EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	_, err = input.OracleKeeper.GetExchangeRate(input.Ctx.WithBlockHeight(1), types.TestDenomD)
	require.Error(t, err)
}

func TestOracleDrop(t *testing.T) {
	input, h := setup(t)

	input.OracleKeeper.SetExchangeRate(input.Ctx, types.TestDenomC, randomExchangeRate)

	// Account 1, DenomC
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate}}, 0)

	// Immediately swap halt after an illiquid oracle vote
	oracle.EndBlocker(input.Ctx, input.OracleKeeper)

	_, err := input.OracleKeeper.GetExchangeRate(input.Ctx, types.TestDenomC)
	require.Error(t, err)
}

func TestOracleTally(t *testing.T) {
	input, _ := setup(t)

	ballot := types.ExchangeRateBallot{}
	rates, valAddrs, stakingKeeper := types.GenerateRandomTestCase()
	input.OracleKeeper.StakingKeeper = stakingKeeper
	h := keeper.NewMsgServerImpl(input.OracleKeeper)
	for i, rate := range rates {

		decExchangeRate := sdk.NewDecWithPrec(int64(rate*math.Pow10(keeper.OracleDecPrecision)), int64(keeper.OracleDecPrecision))
		exchangeRateStr := decExchangeRate.String() + types.TestDenomD

		salt := fmt.Sprintf("%d", i)
		hash := types.GetAggregateVoteHash(salt, exchangeRateStr, valAddrs[i])
		prevoteMsg := types.NewMsgAggregateExchangeRatePrevote(hash, sdk.AccAddress(valAddrs[i]), valAddrs[i])
		voteMsg := types.NewMsgAggregateExchangeRateVote(salt, exchangeRateStr, sdk.AccAddress(valAddrs[i]), valAddrs[i])

		_, err1 := h.AggregateExchangeRatePrevote(input.Ctx.WithBlockHeight(0), prevoteMsg)
		_, err2 := h.AggregateExchangeRateVote(input.Ctx.WithBlockHeight(1), voteMsg)
		require.NoError(t, err1)
		require.NoError(t, err2)

		power := stakingAmt.QuoRaw(types.MicroUnit).Int64()
		if decExchangeRate.IsZero() {
			power = int64(0)
		}

		vote := types.NewVoteForTally(
			decExchangeRate, types.TestDenomD, valAddrs[i], power)
		ballot = append(ballot, vote)

		// change power of every three validator
		if i%3 == 0 {
			stakingKeeper.Validators()[i].SetConsensusPower(int64(i + 1))
		}
	}

	validatorClaimMap := make(map[string]types.Claim)
	for _, valAddr := range valAddrs {
		validatorClaimMap[valAddr.String()] = types.Claim{
			Power:     stakingKeeper.Validator(input.Ctx, valAddr).GetConsensusPower(sdk.DefaultPowerReduction),
			Weight:    int64(0),
			WinCount:  int64(0),
			Recipient: valAddr,
		}
	}
	sort.Sort(ballot)
	weightedMedian, _ := ballot.WeightedMedian()
	standardDeviation, _ := ballot.StandardDeviation()
	maxSpread := weightedMedian.Mul(input.OracleKeeper.RewardBand(input.Ctx).QuoInt64(2))

	if standardDeviation.GT(maxSpread) {
		maxSpread = standardDeviation
	}

	expectedValidatorClaimMap := make(map[string]types.Claim)
	for _, valAddr := range valAddrs {
		expectedValidatorClaimMap[valAddr.String()] = types.Claim{
			Power:     stakingKeeper.Validator(input.Ctx, valAddr).GetConsensusPower(sdk.DefaultPowerReduction),
			Weight:    int64(0),
			WinCount:  int64(0),
			Recipient: valAddr,
		}
	}

	for _, vote := range ballot {
		if (vote.ExchangeRate.GTE(weightedMedian.Sub(maxSpread)) &&
			vote.ExchangeRate.LTE(weightedMedian.Add(maxSpread))) ||
			!vote.ExchangeRate.IsPositive() {
			key := vote.Voter.String()
			claim := expectedValidatorClaimMap[key]
			claim.Weight += vote.Power
			claim.WinCount++
			expectedValidatorClaimMap[key] = claim
		}
	}

	missMap := map[string]sdk.ValAddress{}

	tallyMedian, _ := oracle.Tally(input.Ctx, ballot, input.OracleKeeper.RewardBand(input.Ctx), validatorClaimMap, missMap)

	require.Equal(t, validatorClaimMap, expectedValidatorClaimMap)
	require.Equal(t, tallyMedian.MulInt64(100).TruncateInt(), weightedMedian.MulInt64(100).TruncateInt())
}

func TestOracleTallyTiming(t *testing.T) {
	input, h := setup(t)

	// all the keeper.Addrs vote for the block ... not last period block yet, so tally fails
	for i := range keeper.Addrs[:2] {
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomD, Amount: randomExchangeRate}}, i)
	}

	params := input.OracleKeeper.GetParams(input.Ctx)
	params.VotePeriod = 10 // set vote period to 10 for now, for convenience
	input.OracleKeeper.SetParams(input.Ctx, params)
	require.Equal(t, 0, int(input.Ctx.BlockHeight()))

	oracle.EndBlocker(input.Ctx, input.OracleKeeper)
	_, err := input.OracleKeeper.GetExchangeRate(input.Ctx, types.TestDenomD)
	require.Error(t, err)

	input.Ctx = input.Ctx.WithBlockHeight(int64(params.VotePeriod - 1))

	oracle.EndBlocker(input.Ctx, input.OracleKeeper)
	_, err = input.OracleKeeper.GetExchangeRate(input.Ctx, types.TestDenomD)
	require.NoError(t, err)
}

func TestOracleRewardBand(t *testing.T) {
	input, h := setup(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.Whitelist = types.DenomList{{Name: types.TestDenomC}}
	input.OracleKeeper.SetParams(input.Ctx, params)

	rewardSpread := randomExchangeRate.Mul(input.OracleKeeper.RewardBand(input.Ctx).QuoInt64(2))

	// no one will miss the vote
	// Account 1, DenomC
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate.Sub(rewardSpread)}}, 0)

	// Account 2, DenomC
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate}}, 1)

	// Account 3, DenomC
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate.Add(rewardSpread)}}, 2)

	oracle.EndBlocker(input.Ctx, input.OracleKeeper)

	require.Equal(t, uint64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, uint64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, uint64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))

	// Account 1 will miss the vote due to raward band condition
	// Account 1, DenomC
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate.Sub(rewardSpread.Add(sdk.OneDec()))}}, 0)

	// Account 2, DenomC
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate}}, 1)

	// Account 3, DenomC
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate.Add(rewardSpread)}}, 2)

	oracle.EndBlocker(input.Ctx, input.OracleKeeper)

	require.Equal(t, uint64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, uint64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, uint64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))

	// no one will miss the vote if someone submits an extra denom
	// Account 1, DenomC
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate.Sub(rewardSpread)}}, 0)

	// Account 2, DenomC
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate}}, 1)

	// Account 3, DenomC +  Denom D
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{
		{Denom: types.TestDenomC, Amount: randomExchangeRate.Add(rewardSpread)},
		{Denom: types.TestDenomE, Amount: randomExchangeRate.Add(rewardSpread)},
	}, 2)

	oracle.EndBlocker(input.Ctx, input.OracleKeeper)

	require.Equal(t, uint64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, uint64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, uint64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))

	// no one will miss the vote if there is threshold of an extra denom
	// Account 1, DenomC
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate.Sub(rewardSpread)}}, 0)

	// Account 2, DenomC
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{
		{Denom: types.TestDenomC, Amount: randomExchangeRate},
		{Denom: types.TestDenomE, Amount: randomExchangeRate.Add(rewardSpread)},
	}, 1)

	// Account 3, DenomC +  Denom D
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{
		{Denom: types.TestDenomC, Amount: randomExchangeRate.Add(rewardSpread)},
		{Denom: types.TestDenomE, Amount: randomExchangeRate.Add(rewardSpread)},
	}, 2)

	oracle.EndBlocker(input.Ctx, input.OracleKeeper)

	require.Equal(t, uint64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, uint64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, uint64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))
}

func TestOracleEnsureSorted(t *testing.T) {
	input, h := setup(t)

	for i := 0; i < 100; i++ {
		exchangeRateA1 := sdk.NewDecWithPrec(int64(rand.Uint64()%100000000), 6).MulInt64(types.MicroUnit)
		exchangeRateB1 := sdk.NewDecWithPrec(int64(rand.Uint64()%100000000), 6).MulInt64(types.MicroUnit)

		exchangeRateA2 := sdk.NewDecWithPrec(int64(rand.Uint64()%100000000), 6).MulInt64(types.MicroUnit)
		exchangeRateB2 := sdk.NewDecWithPrec(int64(rand.Uint64()%100000000), 6).MulInt64(types.MicroUnit)

		exchangeRateA3 := sdk.NewDecWithPrec(int64(rand.Uint64()%100000000), 6).MulInt64(types.MicroUnit)
		exchangeRateB3 := sdk.NewDecWithPrec(int64(rand.Uint64()%100000000), 6).MulInt64(types.MicroUnit)

		// Account 1, DenomB, DenomC
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomB, Amount: exchangeRateB1}, {Denom: types.TestDenomC, Amount: exchangeRateA1}}, 0)

		// Account 2, DenomB, DenomC
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomB, Amount: exchangeRateB2}, {Denom: types.TestDenomC, Amount: exchangeRateA2}}, 1)

		// Account 3, DenomB, DenomC
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomB, Amount: exchangeRateA3}, {Denom: types.TestDenomC, Amount: exchangeRateB3}}, 2)

		require.NotPanics(t, func() {
			oracle.EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)
		})
	}
}

func TestInvalidVotesSlashing(t *testing.T) {
	input, h := setup(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.Whitelist = types.DenomList{{Name: types.TestDenomC}}
	input.OracleKeeper.SetParams(input.Ctx, params)

	votePeriodsPerWindow := sdk.NewDec(int64(input.OracleKeeper.SlashWindow(input.Ctx))).QuoInt64(int64(input.OracleKeeper.VotePeriod(input.Ctx))).TruncateInt64()
	slashFraction := input.OracleKeeper.SlashFraction(input.Ctx)
	minValidPerWindow := input.OracleKeeper.MinValidPerWindow(input.Ctx)

	for i := uint64(0); i < uint64(sdk.OneDec().Sub(minValidPerWindow).MulInt64(votePeriodsPerWindow).TruncateInt64()); i++ {
		input.Ctx = input.Ctx.WithBlockHeight(input.Ctx.BlockHeight() + 1)

		// Account 1, DenomC
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate}}, 0)

		// Account 2, DenomC, miss vote
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate.Add(sdk.NewDec(100000000000000))}}, 1)

		// Account 3, DenomC
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate}}, 2)

		oracle.EndBlocker(input.Ctx, input.OracleKeeper)
		require.Equal(t, i+1, input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	}

	validator := input.StakingKeeper.Validator(input.Ctx, keeper.ValAddrs[1])
	require.Equal(t, stakingAmt, validator.GetBondedTokens())

	// one more miss vote will inccur keeper.ValAddrs[1] slashing
	// Account 1, DenomC
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate}}, 0)

	// Account 2, DenomC, miss vote
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate.Add(sdk.NewDec(100000000000000))}}, 1)

	// Account 3, DenomC
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate}}, 2)

	input.Ctx = input.Ctx.WithBlockHeight(votePeriodsPerWindow - 1)
	oracle.EndBlocker(input.Ctx, input.OracleKeeper)
	validator = input.StakingKeeper.Validator(input.Ctx, keeper.ValAddrs[1])
	require.Equal(t, sdk.OneDec().Sub(slashFraction).MulInt(stakingAmt).TruncateInt(), validator.GetBondedTokens())
}

func TestWhitelistSlashing(t *testing.T) {
	input, h := setup(t)

	votePeriodsPerWindow := sdk.NewDec(int64(input.OracleKeeper.SlashWindow(input.Ctx))).QuoInt64(int64(input.OracleKeeper.VotePeriod(input.Ctx))).TruncateInt64()
	slashFraction := input.OracleKeeper.SlashFraction(input.Ctx)
	minValidPerWindow := input.OracleKeeper.MinValidPerWindow(input.Ctx)

	for i := uint64(0); i < uint64(sdk.OneDec().Sub(minValidPerWindow).MulInt64(votePeriodsPerWindow).TruncateInt64()); i++ {
		input.Ctx = input.Ctx.WithBlockHeight(input.Ctx.BlockHeight() + 1)

		// Account 2, DenomC
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate}}, 1)
		// Account 3, DenomC
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate}}, 2)

		oracle.EndBlocker(input.Ctx, input.OracleKeeper)
		require.Equal(t, i+1, input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	}

	validator := input.StakingKeeper.Validator(input.Ctx, keeper.ValAddrs[0])
	require.Equal(t, stakingAmt, validator.GetBondedTokens())

	// one more miss vote will inccur Account 1 slashing

	// Account 2, DenomC
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate}}, 1)
	// Account 3, DenomC
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate}}, 2)

	input.Ctx = input.Ctx.WithBlockHeight(votePeriodsPerWindow - 1)
	oracle.EndBlocker(input.Ctx, input.OracleKeeper)
	validator = input.StakingKeeper.Validator(input.Ctx, keeper.ValAddrs[0])
	require.Equal(t, sdk.OneDec().Sub(slashFraction).MulInt(stakingAmt).TruncateInt(), validator.GetBondedTokens())
}

func TestNotPassedBallotSlashing(t *testing.T) {
	input, h := setup(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.Whitelist = types.DenomList{{Name: types.TestDenomC}}
	input.OracleKeeper.SetParams(input.Ctx, params)

	input.Ctx = input.Ctx.WithBlockHeight(input.Ctx.BlockHeight() + 1)

	// Account 1, DenomC
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate}}, 0)

	oracle.EndBlocker(input.Ctx, input.OracleKeeper)
	// Not slashing accounts that have voted
	require.Equal(t, uint64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, uint64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, uint64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))
}

func TestAbstainSlashing(t *testing.T) {
	input, h := setup(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.Whitelist = types.DenomList{{Name: types.TestDenomC}}
	input.OracleKeeper.SetParams(input.Ctx, params)

	votePeriodsPerWindow := sdk.NewDec(int64(input.OracleKeeper.SlashWindow(input.Ctx))).QuoInt64(int64(input.OracleKeeper.VotePeriod(input.Ctx))).TruncateInt64()
	minValidPerWindow := input.OracleKeeper.MinValidPerWindow(input.Ctx)

	for i := uint64(0); i <= uint64(sdk.OneDec().Sub(minValidPerWindow).MulInt64(votePeriodsPerWindow).TruncateInt64()); i++ {
		input.Ctx = input.Ctx.WithBlockHeight(input.Ctx.BlockHeight() + 1)

		// Account 1, DenomC
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate}}, 0)

		// Account 2, DenomC, abstain vote
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: sdk.ZeroDec()}}, 1)

		// Account 3, DenomC
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate}}, 2)

		oracle.EndBlocker(input.Ctx, input.OracleKeeper)
		require.Equal(t, uint64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	}

	validator := input.StakingKeeper.Validator(input.Ctx, keeper.ValAddrs[1])
	require.Equal(t, stakingAmt, validator.GetBondedTokens())
}

func TestVoteTargets(t *testing.T) {
	input, h := setup(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.Whitelist = types.DenomList{{Name: types.TestDenomC}, {Name: types.TestDenomD}}
	input.OracleKeeper.SetParams(input.Ctx, params)

	// DenomC
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate}}, 0)
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate}}, 1)
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate}}, 2)

	oracle.EndBlocker(input.Ctx, input.OracleKeeper)

	// missed D
	require.Equal(t, uint64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, uint64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, uint64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))

	// delete DenomD
	params.Whitelist = types.DenomList{{Name: types.TestDenomC}}
	input.OracleKeeper.SetParams(input.Ctx, params)

	// DenomC, missing
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{}, 0)
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{}, 1)
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{}, 2)

	oracle.EndBlocker(input.Ctx, input.OracleKeeper)

	require.Equal(t, uint64(2), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, uint64(2), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, uint64(2), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))

	// DenomC, no missing
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate}}, 0)
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate}}, 1)
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: randomExchangeRate}}, 2)

	oracle.EndBlocker(input.Ctx, input.OracleKeeper)

	require.Equal(t, uint64(2), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, uint64(2), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, uint64(2), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))
}

func TestAbstainWithSmallStakingPower(t *testing.T) {
	input, h := setupWithSmallVotingPower(t)

	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: types.TestDenomC, Amount: sdk.ZeroDec()}}, 0)

	oracle.EndBlocker(input.Ctx, input.OracleKeeper)
	_, err := input.OracleKeeper.GetExchangeRate(input.Ctx, types.TestDenomC)
	require.Error(t, err)
}

func makeAggregatePrevoteAndVote(t *testing.T, input keeper.TestInput, h types.MsgServer, height int64, rates sdk.DecCoins, idx int) {
	// Account 1, DenomD
	salt := "fc5bb0bc63e54b2918d9334bf3259f5dc575e8d7a4df4e836dd80f1ad62aa89b"
	hash := types.GetAggregateVoteHash(salt, rates.String(), keeper.ValAddrs[idx])

	prevoteMsg := types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[idx], keeper.ValAddrs[idx])
	_, err := h.AggregateExchangeRatePrevote(input.Ctx.WithBlockHeight(height), prevoteMsg)
	require.NoError(t, err)

	voteMsg := types.NewMsgAggregateExchangeRateVote(salt, rates.String(), keeper.Addrs[idx], keeper.ValAddrs[idx])
	_, err = h.AggregateExchangeRateVote(input.Ctx.WithBlockHeight(height+1), voteMsg)
	require.NoError(t, err)
}
