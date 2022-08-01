package wasm_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

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

	querier := wasm.NewWasmQuerier(input.OracleKeeper)
	var err error

	// empty data will occur error
	_, err = querier.QueryCustom(input.Ctx, []byte{})
	require.Error(t, err)

	// not existing quote denom query
	queryParams := wasm.ExchangeRateQueryParams{
		BaseDenom:   types.TestDenomA,
		QuoteDenoms: []string{types.TestDenomI},
	}
	bz, err := json.Marshal(wasm.CosmosQuery{
		ExchangeRates: &queryParams,
	})
	require.NoError(t, err)

	res, err := querier.QueryCustom(input.Ctx, bz)
	require.NoError(t, err)

	var exchangeRatesResponse wasm.ExchangeRatesQueryResponse
	err = json.Unmarshal(res, &exchangeRatesResponse)
	require.NoError(t, err)
	require.Equal(t, wasm.ExchangeRatesQueryResponse{
		BaseDenom:     types.TestDenomA,
		ExchangeRates: nil,
	}, exchangeRatesResponse)

	// not existing base denom query
	queryParams = wasm.ExchangeRateQueryParams{
		BaseDenom:   types.TestDenomE,
		QuoteDenoms: []string{types.TestDenomC, types.TestDenomB, types.TestDenomD},
	}
	bz, err = json.Marshal(wasm.CosmosQuery{
		ExchangeRates: &queryParams,
	})
	require.NoError(t, err)

	res, err = querier.QueryCustom(input.Ctx, bz)
	require.Error(t, err)

	// valid query luna exchange rates
	queryParams = wasm.ExchangeRateQueryParams{
		BaseDenom:   types.TestDenomA,
		QuoteDenoms: []string{types.TestDenomC, types.TestDenomB, types.TestDenomD},
	}
	bz, err = json.Marshal(wasm.CosmosQuery{
		ExchangeRates: &queryParams,
	})
	require.NoError(t, err)

	res, err = querier.QueryCustom(input.Ctx, bz)
	require.NoError(t, err)

	err = json.Unmarshal(res, &exchangeRatesResponse)
	require.NoError(t, err)
	require.Equal(t, exchangeRatesResponse, wasm.ExchangeRatesQueryResponse{
		BaseDenom: types.TestDenomA,
		ExchangeRates: []wasm.ExchangeRateItem{
			{
				ExchangeRate: ExchangeRateC.String(),
				QuoteDenom:   types.TestDenomC,
			},
			{
				ExchangeRate: ExchangeRateB.String(),
				QuoteDenom:   types.TestDenomB,
			},
			{
				ExchangeRate: ExchangeRateD.String(),
				QuoteDenom:   types.TestDenomD,
			},
		},
	})
}
