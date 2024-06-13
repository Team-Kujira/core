package keeper

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/Team-Kujira/core/x/oracle/types"
	"github.com/stretchr/testify/require"
)

func TestHistoricalRates(t *testing.T) {
	input := CreateTestInput(t)

	now := time.Now()
	historicalRates := []types.HistoricalExchangeRate{
		{
			Epoch:        "day",
			Timestamp:    now.Unix() - 86400,
			Denom:        "BTC",
			ExchangeRate: math.LegacyNewDec(64000),
		},
		{
			Epoch:        "day",
			Timestamp:    now.Unix(),
			Denom:        "BTC",
			ExchangeRate: math.LegacyNewDec(63588),
		},
		{
			Epoch:        "hour",
			Timestamp:    now.Unix() - 3600,
			Denom:        "BTC",
			ExchangeRate: math.LegacyNewDec(63800),
		},
		{
			Epoch:        "hour",
			Timestamp:    now.Unix(),
			Denom:        "BTC",
			ExchangeRate: math.LegacyNewDec(63588),
		},
	}
	for _, rate := range historicalRates {
		input.OracleKeeper.SetHistoricalExchangeRate(input.Ctx, rate)
	}

	for _, rate := range historicalRates {
		storedRate, err := input.OracleKeeper.GetHistoricalExchangeRate(input.Ctx, rate.Epoch, rate.Denom, rate.Timestamp)
		require.NoError(t, err)
		require.Equal(t, storedRate.ExchangeRate.String(), rate.ExchangeRate.String())
	}

	rate := input.OracleKeeper.OldestHistoricalExchangeRateByEpochDenom(input.Ctx, "day", "BTC")
	require.Equal(t, rate.ExchangeRate.String(), historicalRates[0].ExchangeRate.String())
	rate = input.OracleKeeper.LatestHistoricalExchangeRateByEpochDenom(input.Ctx, "day", "BTC")
	require.Equal(t, rate.ExchangeRate.String(), historicalRates[1].ExchangeRate.String())

	input.OracleKeeper.DeleteHistoricalExchangeRate(input.Ctx, historicalRates[0])
	rate = input.OracleKeeper.OldestHistoricalExchangeRateByEpochDenom(input.Ctx, "day", "BTC")
	require.Equal(t, rate.ExchangeRate.String(), historicalRates[1].ExchangeRate.String())
	rate = input.OracleKeeper.LatestHistoricalExchangeRateByEpochDenom(input.Ctx, "day", "BTC")
	require.Equal(t, rate.ExchangeRate.String(), historicalRates[1].ExchangeRate.String())
}
