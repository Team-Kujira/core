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

	KRWExchangeRate := sdk.NewDec(1700)
	USDExchangeRate := sdk.NewDecWithPrec(17, 1)
	SDRExchangeRate := sdk.NewDecWithPrec(19, 1)
	input.OracleKeeper.SetExchangeRate(input.Ctx, types.TestDenomC, KRWExchangeRate)
	input.OracleKeeper.SetExchangeRate(input.Ctx, types.TestDenomB, USDExchangeRate)
	input.OracleKeeper.SetExchangeRate(input.Ctx, types.TestDenomD, SDRExchangeRate)

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
				ExchangeRate: KRWExchangeRate.String(),
				QuoteDenom:   types.TestDenomC,
			},
			{
				ExchangeRate: USDExchangeRate.String(),
				QuoteDenom:   types.TestDenomB,
			},
			{
				ExchangeRate: SDRExchangeRate.String(),
				QuoteDenom:   types.TestDenomD,
			},
		},
	})
}
