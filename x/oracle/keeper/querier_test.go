package keeper

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Team-Kujira/core/x/oracle/types"
)

func TestQueryParams(t *testing.T) {
	input := CreateTestInput(t)

	querier := NewQuerier(input.OracleKeeper)
	res, err := querier.Params(input.Ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)

	require.Equal(t, input.OracleKeeper.GetParams(input.Ctx), res.Params)
}

func TestQueryExchangeRate(t *testing.T) {
	input := CreateTestInput(t)
	querier := NewQuerier(input.OracleKeeper)

	rate := math.LegacyNewDec(1700)
	input.OracleKeeper.SetExchangeRate(input.Ctx, types.TestDenomD, rate)

	// empty request
	_, err := querier.ExchangeRate(input.Ctx, nil)
	require.Error(t, err)

	// Query to grpc
	res, err := querier.ExchangeRate(input.Ctx, &types.QueryExchangeRateRequest{
		Denom: types.TestDenomD,
	})
	require.NoError(t, err)
	require.Equal(t, rate, res.ExchangeRate)
}

func TestQueryMissCounter(t *testing.T) {
	input := CreateTestInput(t)
	querier := NewQuerier(input.OracleKeeper)

	missCounter := uint64(1)
	input.OracleKeeper.SetMissCounter(input.Ctx, ValAddrs[0], missCounter)

	// empty request
	_, err := querier.MissCounter(input.Ctx, nil)
	require.Error(t, err)

	// Query to grpc
	res, err := querier.MissCounter(input.Ctx, &types.QueryMissCounterRequest{
		ValidatorAddr: ValAddrs[0].String(),
	})
	require.NoError(t, err)
	require.Equal(t, missCounter, res.MissCounter)
}

func TestQueryExchangeRates(t *testing.T) {
	input := CreateTestInput(t)
	querier := NewQuerier(input.OracleKeeper)

	rate := math.LegacyNewDec(1700)
	input.OracleKeeper.SetExchangeRate(input.Ctx, types.TestDenomD, rate)
	input.OracleKeeper.SetExchangeRate(input.Ctx, types.TestDenomB, rate)

	res, err := querier.ExchangeRates(input.Ctx, &types.QueryExchangeRatesRequest{})
	require.NoError(t, err)

	require.Equal(t, sdk.DecCoins{
		sdk.NewDecCoinFromDec(types.TestDenomB, rate),
		sdk.NewDecCoinFromDec(types.TestDenomD, rate),
	}, res.ExchangeRates)
}

func TestQueryActives(t *testing.T) {
	input := CreateTestInput(t)
	querier := NewQuerier(input.OracleKeeper)

	rate := math.LegacyNewDec(1700)
	input.OracleKeeper.SetExchangeRate(input.Ctx, types.TestDenomD, rate)
	input.OracleKeeper.SetExchangeRate(input.Ctx, types.TestDenomC, rate)
	input.OracleKeeper.SetExchangeRate(input.Ctx, types.TestDenomB, rate)

	res, err := querier.Actives(input.Ctx, &types.QueryActivesRequest{})
	require.NoError(t, err)

	targetDenoms := []string{
		types.TestDenomB,
		types.TestDenomC,
		types.TestDenomD,
	}

	require.Equal(t, targetDenoms, res.Actives)
}
