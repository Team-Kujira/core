package wasmbinding_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"

	"github.com/Team-Kujira/core/x/oracle/types"
	"github.com/Team-Kujira/core/x/oracle/wasm"

	"github.com/Team-Kujira/core/app"
	"github.com/Team-Kujira/core/wasmbinding"
	"github.com/Team-Kujira/core/wasmbinding/bindings"

	tmtypes "github.com/cometbft/cometbft/proto/tendermint/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func TestQueryExchangeRates(t *testing.T) {
	app := app.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmtypes.Header{Height: 1, ChainID: "kujira-1", Time: time.Now().UTC()})

	ExchangeRateC := sdk.NewDec(1700)
	ExchangeRateB := sdk.NewDecWithPrec(17, 1)
	ExchangeRateD := sdk.NewDecWithPrec(19, 1)
	app.OracleKeeper.SetExchangeRate(ctx, types.TestDenomA, sdk.NewDec(1))
	app.OracleKeeper.SetExchangeRate(ctx, types.TestDenomC, ExchangeRateC)
	app.OracleKeeper.SetExchangeRate(ctx, types.TestDenomB, ExchangeRateB)
	app.OracleKeeper.SetExchangeRate(ctx, types.TestDenomD, ExchangeRateD)

	plugin := wasmbinding.NewQueryPlugin(app.BankKeeper, app.OracleKeeper, *app.DenomKeeper, *app.IBCKeeper, app.InterTxKeeper, app.GetKey(ibcexported.StoreKey))
	querier := wasmbinding.CustomQuerier(plugin)
	var err error

	// empty data will occur error
	_, err = querier(ctx, []byte{})
	require.Error(t, err)

	// not existing quote denom query
	queryParams := wasm.ExchangeRateQueryParams{
		Denom: types.TestDenomI,
	}
	bz, err := json.Marshal(bindings.CosmosQuery{
		Oracle: &wasm.OracleQuery{
			ExchangeRate: &queryParams,
		},
	})
	require.NoError(t, err)

	res, err := querier(ctx, bz)
	require.Error(t, err)

	var exchangeRatesResponse wasm.ExchangeRateQueryResponse
	err = json.Unmarshal(res, &exchangeRatesResponse)
	require.Error(t, err)

	// not existing base denom query
	queryParams = wasm.ExchangeRateQueryParams{
		Denom: types.TestDenomC,
	}
	bz, err = json.Marshal(bindings.CosmosQuery{
		Oracle: &wasm.OracleQuery{
			ExchangeRate: &queryParams,
		},
	})
	require.NoError(t, err)

	res, err = querier(ctx, bz)
	require.NoError(t, err)

	queryParams = wasm.ExchangeRateQueryParams{
		Denom: types.TestDenomB,
	}
	bz, err = json.Marshal(bindings.CosmosQuery{
		Oracle: &wasm.OracleQuery{
			ExchangeRate: &queryParams,
		},
	})
	require.NoError(t, err)

	res, err = querier(ctx, bz)
	require.NoError(t, err)

	err = json.Unmarshal(res, &exchangeRatesResponse)
	require.NoError(t, err)
	require.Equal(t, exchangeRatesResponse, wasm.ExchangeRateQueryResponse{
		Rate: ExchangeRateB.String(),
	})
}

func TestSupply(t *testing.T) {
	app := app.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmtypes.Header{Height: 1, ChainID: "kujira-1", Time: time.Now().UTC()})

	plugin := wasmbinding.NewQueryPlugin(app.BankKeeper, app.OracleKeeper, *app.DenomKeeper, *app.IBCKeeper, app.InterTxKeeper, app.GetKey(ibcexported.StoreKey))
	querier := wasmbinding.CustomQuerier(plugin)

	var err error

	// empty data will occur error
	_, err = querier(ctx, []byte{})
	require.Error(t, err)

	queryParams := banktypes.QuerySupplyOfRequest{
		Denom: types.TestDenomA,
	}
	bz, err := json.Marshal(bindings.CosmosQuery{
		Bank: &bindings.BankQuery{
			Supply: &queryParams,
		},
	})
	require.NoError(t, err)
	var x banktypes.QuerySupplyOfResponse

	res, err := querier(ctx, bz)

	err = json.Unmarshal(res, &x)
	require.NoError(t, err)
}

func TestVerifyMembership(t *testing.T) {
	app := app.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmtypes.Header{Height: 1, ChainID: "kujira-1", Time: time.Now().UTC()})

	var err error
	contractAddr, _ := SetupContract(t, ctx, app)
	queryMsg := bindings.IbcQuery{
		VerifyMembership: &bindings.VerifyMembershipQuery{
			Connection:     "connection-0",
			RevisionNumber: 0,
			RevisionHeight: 0,
			Proof:          []byte{},
			Value:          []byte{},
			PathPrefix:     "ibc",
			PathKey:        "connections/connection-0",
		},
	}

	bz, err := json.Marshal(queryMsg)
	require.NoError(t, err)

	bz, err = app.WasmKeeper.QuerySmart(ctx, contractAddr, bz)
	require.Error(t, err)
}

func TestVerifyNonMembership(t *testing.T) {
	app := app.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmtypes.Header{Height: 1, ChainID: "kujira-1", Time: time.Now().UTC()})

	var err error
	contractAddr, _ := SetupContract(t, ctx, app)
	queryMsg := bindings.IbcQuery{
		VerifyNonMembership: &bindings.VerifyNonMembershipQuery{
			Connection:     "connection-0",
			RevisionNumber: 0,
			RevisionHeight: 0,
			Proof:          []byte{},
			PathPrefix:     "ibc",
			PathKey:        "connections/connection-0",
		},
	}

	bz, err := json.Marshal(queryMsg)
	require.NoError(t, err)

	bz, err = app.WasmKeeper.QuerySmart(ctx, contractAddr, bz)
	require.Error(t, err)
}
