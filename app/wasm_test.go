package app_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"kujira/app"
	"kujira/x/oracle/keeper"
	"kujira/x/oracle/types"
	"kujira/x/oracle/wasm"
)

func TestQueryExchangeRates(t *testing.T) {
	input := keeper.CreateTestInput(t)

	ExchangeRateC := sdk.NewDec(1700)
	ExchangeRateB := sdk.NewDecWithPrec(17, 1)
	ExchangeRateD := sdk.NewDecWithPrec(19, 1)
	input.OracleKeeper.SetExchangeRate(input.Ctx, types.TestDenomA, sdk.NewDec(1))
	input.OracleKeeper.SetExchangeRate(input.Ctx, types.TestDenomC, ExchangeRateC)
	input.OracleKeeper.SetExchangeRate(input.Ctx, types.TestDenomB, ExchangeRateB)
	input.OracleKeeper.SetExchangeRate(input.Ctx, types.TestDenomD, ExchangeRateD)

	querier := app.NewWasmQuerier(input.OracleKeeper)
	var err error

	// empty data will occur error
	_, err = querier.QueryCustom(input.Ctx, []byte{})
	require.Error(t, err)

	// not existing quote denom query
	queryParams := wasm.ExchangeRateQueryParams{
		Denom: types.TestDenomI,
	}
	bz, err := json.Marshal(wasm.CosmosQuery{
		Oracle: &wasm.OracleQuery{
			ExchangeRate: &queryParams,
		},
	})
	require.NoError(t, err)

	res, err := querier.QueryCustom(input.Ctx, bz)
	require.Error(t, err)

	var exchangeRatesResponse wasm.ExchangeRateQueryResponse
	err = json.Unmarshal(res, &exchangeRatesResponse)
	require.Error(t, err)

	// not existing base denom query
	queryParams = wasm.ExchangeRateQueryParams{
		Denom: types.TestDenomC,
	}
	bz, err = json.Marshal(wasm.CosmosQuery{
		Oracle: &wasm.OracleQuery{
			ExchangeRate: &queryParams,
		},
	})
	require.NoError(t, err)

	res, err = querier.QueryCustom(input.Ctx, bz)
	require.NoError(t, err)

	queryParams = wasm.ExchangeRateQueryParams{
		Denom: types.TestDenomB,
	}
	bz, err = json.Marshal(wasm.CosmosQuery{
		Oracle: &wasm.OracleQuery{
			ExchangeRate: &queryParams,
		},
	})
	require.NoError(t, err)

	res, err = querier.QueryCustom(input.Ctx, bz)
	require.NoError(t, err)

	err = json.Unmarshal(res, &exchangeRatesResponse)
	require.NoError(t, err)
	require.Equal(t, exchangeRatesResponse, wasm.ExchangeRateQueryResponse{
		Rate: ExchangeRateB.String(),
	})
}
