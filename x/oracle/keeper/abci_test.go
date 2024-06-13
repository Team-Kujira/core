package keeper

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/Team-Kujira/core/x/oracle/types"
	"github.com/stretchr/testify/require"
)

func TestEndBlocker(t *testing.T) {
	input := CreateTestInput(t)

	now := time.Now()
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.RequiredDenoms = []string{"BTC"}
	params.ExchangeRateSnapEpochs = []types.PriceIntervalParam{
		{
			Epoch:    "day",
			Duration: 86400,
			MaxCount: 3,
		},
		{
			Epoch:    "hour",
			Duration: 3600,
			MaxCount: 3,
		},
	}
	input.OracleKeeper.SetParams(input.Ctx, params)

	// Run endblocker now
	input.Ctx = input.Ctx.WithBlockTime(now)
	input.OracleKeeper.SetExchangeRate(input.Ctx, "BTC", math.LegacyNewDec(63500))
	input.OracleKeeper.EndBlocker(input.Ctx)
	rate := input.OracleKeeper.LatestHistoricalExchangeRateByEpochDenom(input.Ctx, "day", "BTC")
	require.Equal(t, rate.ExchangeRate.String(), math.LegacyNewDec(63500).String())
	rate = input.OracleKeeper.OldestHistoricalExchangeRateByEpochDenom(input.Ctx, "day", "BTC")
	require.Equal(t, rate.ExchangeRate.String(), math.LegacyNewDec(63500).String())
	rate = input.OracleKeeper.LatestHistoricalExchangeRateByEpochDenom(input.Ctx, "hour", "BTC")
	require.Equal(t, rate.ExchangeRate.String(), math.LegacyNewDec(63500).String())
	rate = input.OracleKeeper.OldestHistoricalExchangeRateByEpochDenom(input.Ctx, "hour", "BTC")
	require.Equal(t, rate.ExchangeRate.String(), math.LegacyNewDec(63500).String())

	// Price changes after 1 hour
	input.Ctx = input.Ctx.WithBlockTime(now.Add(time.Hour))
	input.OracleKeeper.SetExchangeRate(input.Ctx, "BTC", math.LegacyNewDec(63800))
	input.OracleKeeper.EndBlocker(input.Ctx)
	rate = input.OracleKeeper.LatestHistoricalExchangeRateByEpochDenom(input.Ctx, "day", "BTC")
	require.Equal(t, rate.ExchangeRate.String(), math.LegacyNewDec(63500).String())
	rate = input.OracleKeeper.OldestHistoricalExchangeRateByEpochDenom(input.Ctx, "day", "BTC")
	require.Equal(t, rate.ExchangeRate.String(), math.LegacyNewDec(63500).String())
	rate = input.OracleKeeper.LatestHistoricalExchangeRateByEpochDenom(input.Ctx, "hour", "BTC")
	require.Equal(t, rate.ExchangeRate.String(), math.LegacyNewDec(63800).String())
	rate = input.OracleKeeper.OldestHistoricalExchangeRateByEpochDenom(input.Ctx, "hour", "BTC")
	require.Equal(t, rate.ExchangeRate.String(), math.LegacyNewDec(63500).String())

	// Price changes after 2 hours
	input.Ctx = input.Ctx.WithBlockTime(now.Add(time.Hour * 2))
	input.OracleKeeper.SetExchangeRate(input.Ctx, "BTC", math.LegacyNewDec(64000))
	input.OracleKeeper.EndBlocker(input.Ctx)
	rate = input.OracleKeeper.LatestHistoricalExchangeRateByEpochDenom(input.Ctx, "day", "BTC")
	require.Equal(t, rate.ExchangeRate.String(), math.LegacyNewDec(63500).String())
	rate = input.OracleKeeper.OldestHistoricalExchangeRateByEpochDenom(input.Ctx, "day", "BTC")
	require.Equal(t, rate.ExchangeRate.String(), math.LegacyNewDec(63500).String())
	rate = input.OracleKeeper.LatestHistoricalExchangeRateByEpochDenom(input.Ctx, "hour", "BTC")
	require.Equal(t, rate.ExchangeRate.String(), math.LegacyNewDec(64000).String())
	rate = input.OracleKeeper.OldestHistoricalExchangeRateByEpochDenom(input.Ctx, "hour", "BTC")
	require.Equal(t, rate.ExchangeRate.String(), math.LegacyNewDec(63500).String())

	// Price changes after 1 day
	input.Ctx = input.Ctx.WithBlockTime(now.Add(time.Hour * 24))
	input.OracleKeeper.SetExchangeRate(input.Ctx, "BTC", math.LegacyNewDec(64500))
	input.OracleKeeper.EndBlocker(input.Ctx)
	rate = input.OracleKeeper.LatestHistoricalExchangeRateByEpochDenom(input.Ctx, "day", "BTC")
	require.Equal(t, rate.ExchangeRate.String(), math.LegacyNewDec(64500).String())
	rate = input.OracleKeeper.OldestHistoricalExchangeRateByEpochDenom(input.Ctx, "day", "BTC")
	require.Equal(t, rate.ExchangeRate.String(), math.LegacyNewDec(63500).String())
	rate = input.OracleKeeper.LatestHistoricalExchangeRateByEpochDenom(input.Ctx, "hour", "BTC")
	require.Equal(t, rate.ExchangeRate.String(), math.LegacyNewDec(64500).String())
	rate = input.OracleKeeper.OldestHistoricalExchangeRateByEpochDenom(input.Ctx, "hour", "BTC")
	require.Equal(t, rate.ExchangeRate.String(), math.LegacyNewDec(63800).String())
}
