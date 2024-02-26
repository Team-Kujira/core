package keeper

import (
	"math"
	"sort"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/Team-Kujira/core/x/oracle/types"
)

func TestOracleTally(t *testing.T) {
	input, _ := setup(t)

	ballot := types.ExchangeRateBallot{}
	rates, valAddrs, stakingKeeper := types.GenerateRandomTestCase()
	input.OracleKeeper.StakingKeeper = stakingKeeper
	for i, rate := range rates {
		decExchangeRate := sdkmath.LegacyNewDecWithPrec(int64(rate*math.Pow10(OracleDecPrecision)), int64(OracleDecPrecision))

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
		val, _ := stakingKeeper.Validator(input.Ctx, valAddr)
		validatorClaimMap[valAddr.String()] = types.Claim{
			Power:     val.GetConsensusPower(sdk.DefaultPowerReduction),
			Weight:    int64(0),
			WinCount:  int64(0),
			Recipient: valAddr,
		}
	}
	sort.Sort(ballot)
	weightedMedian, _ := ballot.WeightedMedian()
	standardDeviation, _ := ballot.StandardDeviation()
	maxSpread := weightedMedian.Mul(input.OracleKeeper.MaxDeviation(input.Ctx))

	if standardDeviation.GT(maxSpread) {
		maxSpread = standardDeviation
	}

	expectedValidatorClaimMap := make(map[string]types.Claim)
	for _, valAddr := range valAddrs {
		val, _ := stakingKeeper.Validator(input.Ctx, valAddr)
		expectedValidatorClaimMap[valAddr.String()] = types.Claim{
			Power:     val.GetConsensusPower(sdk.DefaultPowerReduction),
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

	tallyMedian, _ := Tally(input.Ctx, ballot, input.OracleKeeper.MaxDeviation(input.Ctx), validatorClaimMap, missMap)

	require.Equal(t, validatorClaimMap, expectedValidatorClaimMap)
	require.Equal(t, tallyMedian.MulInt64(100).TruncateInt(), weightedMedian.MulInt64(100).TruncateInt())
}
