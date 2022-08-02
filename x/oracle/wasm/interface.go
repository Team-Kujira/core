package wasm

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	"kujira/x/oracle/keeper"
)

// WasmQuerier - staking query interface for wasm contract
type WasmQuerier struct {
	keeper keeper.Keeper
}

// NewWasmQuerier return bank wasm query interface
func NewWasmQuerier(keeper keeper.Keeper) WasmQuerier {
	return WasmQuerier{keeper}
}

// Query - implement query function
func (WasmQuerier) Query(_ sdk.Context, _ wasmvmtypes.QueryRequest) ([]byte, error) {
	return nil, nil
}

// ExchangeRateQueryParams query request params for exchange rates
type ExchangeRateQueryParams struct {
	Denom string `json:"denom"`
}

type CosmosQuery struct {
	Oracle *OracleQuery
}

// OracleQuery custom query interface for oracle querier
type OracleQuery struct {
	ExchangeRate *ExchangeRateQueryParams `json:"exchange_rate,omitempty"`
}

// ExchangeRateQueryResponse - exchange rates query response item
type ExchangeRateQueryResponse struct {
	Rate string `json:"rate"`
}

// QueryCustom implements custom query interface
func (querier WasmQuerier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	var params CosmosQuery
	err := json.Unmarshal(data, &params)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	if params.Oracle.ExchangeRate != nil {
		rate, err := querier.keeper.GetExchangeRate(ctx, params.Oracle.ExchangeRate.Denom)
		if err != nil {
			return nil, err
		}

		bz, err := json.Marshal(ExchangeRateQueryResponse{
			Rate: rate.String(),
		})

		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}

		return bz, nil
	}

	return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown Oracle variant"}
}
