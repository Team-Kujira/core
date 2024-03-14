package keeper

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/Team-Kujira/core/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestExchangeRate(t *testing.T) {
	input := CreateTestInput(t)

	exchangeRateE := math.LegacyNewDecWithPrec(839, int64(OracleDecPrecision)).MulInt64(types.MicroUnit)
	exchangeRateH := math.LegacyNewDecWithPrec(4995, int64(OracleDecPrecision)).MulInt64(types.MicroUnit)
	exchangeRateC := math.LegacyNewDecWithPrec(2838, int64(OracleDecPrecision)).MulInt64(types.MicroUnit)
	exchangeRateA := math.LegacyNewDecWithPrec(3282384, int64(OracleDecPrecision)).MulInt64(types.MicroUnit)

	// Set & get rates
	input.OracleKeeper.SetExchangeRate(input.Ctx, types.TestDenomE, exchangeRateE)
	rate, err := input.OracleKeeper.GetExchangeRate(input.Ctx, types.TestDenomE)
	require.NoError(t, err)
	require.Equal(t, exchangeRateE, rate)

	input.OracleKeeper.SetExchangeRate(input.Ctx, types.TestDenomH, exchangeRateH)
	rate, err = input.OracleKeeper.GetExchangeRate(input.Ctx, types.TestDenomH)
	require.NoError(t, err)
	require.Equal(t, exchangeRateH, rate)

	input.OracleKeeper.SetExchangeRate(input.Ctx, types.TestDenomC, exchangeRateC)
	rate, err = input.OracleKeeper.GetExchangeRate(input.Ctx, types.TestDenomC)
	require.NoError(t, err)
	require.Equal(t, exchangeRateC, rate)

	input.OracleKeeper.SetExchangeRate(input.Ctx, types.TestDenomA, exchangeRateA)
	rate, _ = input.OracleKeeper.GetExchangeRate(input.Ctx, types.TestDenomA)
	require.Equal(t, exchangeRateA, rate)

	input.OracleKeeper.DeleteExchangeRate(input.Ctx, types.TestDenomC)
	_, err = input.OracleKeeper.GetExchangeRate(input.Ctx, types.TestDenomC)
	require.Error(t, err)

	numExchangeRates := 0
	handler := func(denom string, exchangeRate math.LegacyDec) (stop bool) {
		numExchangeRates = numExchangeRates + 1
		return false
	}
	input.OracleKeeper.IterateExchangeRates(input.Ctx, handler)

	require.True(t, numExchangeRates == 3)
}

func TestIterateExchangeRates(t *testing.T) {
	input := CreateTestInput(t)

	exchangeRateE := math.LegacyNewDecWithPrec(839, int64(OracleDecPrecision)).MulInt64(types.MicroUnit)
	exchangeRateH := math.LegacyNewDecWithPrec(4995, int64(OracleDecPrecision)).MulInt64(types.MicroUnit)
	exchangeRateC := math.LegacyNewDecWithPrec(2838, int64(OracleDecPrecision)).MulInt64(types.MicroUnit)
	exchangeRateA := math.LegacyNewDecWithPrec(3282384, int64(OracleDecPrecision)).MulInt64(types.MicroUnit)

	// Set & get rates
	input.OracleKeeper.SetExchangeRate(input.Ctx, types.TestDenomE, exchangeRateE)
	input.OracleKeeper.SetExchangeRate(input.Ctx, types.TestDenomH, exchangeRateH)
	input.OracleKeeper.SetExchangeRate(input.Ctx, types.TestDenomC, exchangeRateC)
	input.OracleKeeper.SetExchangeRate(input.Ctx, types.TestDenomA, exchangeRateA)

	input.OracleKeeper.IterateExchangeRates(input.Ctx, func(denom string, rate math.LegacyDec) (stop bool) {
		switch denom {
		case types.TestDenomE:
			require.Equal(t, exchangeRateE, rate)
		case types.TestDenomH:
			require.Equal(t, exchangeRateH, rate)
		case types.TestDenomC:
			require.Equal(t, exchangeRateC, rate)
		case types.TestDenomA:
			require.Equal(t, exchangeRateA, rate)
		}
		return false
	})
}

func TestRewardPool(t *testing.T) {
	input := CreateTestInput(t)

	fees := sdk.NewCoins(sdk.NewCoin(types.TestDenomD, math.NewInt(1000)))
	acc := input.AccountKeeper.GetModuleAccount(input.Ctx, types.ModuleName)
	err := FundAccount(input, acc.GetAddress(), fees)
	if err != nil {
		panic(err) // never occurs
	}

	KFees := input.OracleKeeper.GetRewardPool(input.Ctx, types.TestDenomD)
	require.Equal(t, fees[0], KFees)
}

func TestParams(t *testing.T) {
	input := CreateTestInput(t)

	// Test default params setting
	input.OracleKeeper.SetParams(input.Ctx, types.DefaultParams())
	params := input.OracleKeeper.GetParams(input.Ctx)
	require.NotNil(t, params)

	// Test custom params setting
	votePeriod := uint64(10)
	voteThreshold := math.LegacyNewDecWithPrec(70, 2)
	maxDeviation := math.LegacyNewDecWithPrec(1, 1)
	slashFraction := math.LegacyNewDecWithPrec(1, 2)
	slashWindow := uint64(1000)
	minValidPerWindow := math.LegacyNewDecWithPrec(1, 4)
	requiredDenoms := []string{
		types.TestDenomD,
		types.TestDenomC,
	}

	// Should really test validateParams, but skipping because obvious
	newParams := types.Params{
		VotePeriod:        votePeriod,
		VoteThreshold:     voteThreshold,
		MaxDeviation:      maxDeviation,
		RequiredDenoms:    requiredDenoms,
		SlashFraction:     slashFraction,
		SlashWindow:       slashWindow,
		MinValidPerWindow: minValidPerWindow,
		Whitelist:         nil,
		RewardBand:        math.LegacyZeroDec(),
	}
	err := input.OracleKeeper.SetParams(input.Ctx, newParams)
	require.NoError(t, err)
	storedParams := input.OracleKeeper.GetParams(input.Ctx)
	require.NotNil(t, storedParams)
	require.Equal(t, storedParams, newParams)
}

func TestMissCounter(t *testing.T) {
	input := CreateTestInput(t)

	// Test default getters and setters
	counter := input.OracleKeeper.GetMissCounter(input.Ctx, ValAddrs[0])
	require.Equal(t, uint64(0), counter)

	missCounter := uint64(10)
	input.OracleKeeper.SetMissCounter(input.Ctx, ValAddrs[0], missCounter)
	counter = input.OracleKeeper.GetMissCounter(input.Ctx, ValAddrs[0])
	require.Equal(t, missCounter, counter)

	input.OracleKeeper.DeleteMissCounter(input.Ctx, ValAddrs[0])
	counter = input.OracleKeeper.GetMissCounter(input.Ctx, ValAddrs[0])
	require.Equal(t, uint64(0), counter)
}

func TestIterateMissCounters(t *testing.T) {
	input := CreateTestInput(t)

	// Test default getters and setters
	counter := input.OracleKeeper.GetMissCounter(input.Ctx, ValAddrs[0])
	require.Equal(t, uint64(0), counter)

	missCounter := uint64(10)
	input.OracleKeeper.SetMissCounter(input.Ctx, ValAddrs[1], missCounter)

	var operators []sdk.ValAddress
	var missCounters []uint64
	input.OracleKeeper.IterateMissCounters(input.Ctx, func(delegator sdk.ValAddress, missCounter uint64) (stop bool) {
		operators = append(operators, delegator)
		missCounters = append(missCounters, missCounter)
		return false
	})

	require.Equal(t, 1, len(operators))
	require.Equal(t, 1, len(missCounters))
	require.Equal(t, ValAddrs[1], operators[0])
	require.Equal(t, missCounter, missCounters[0])
}
