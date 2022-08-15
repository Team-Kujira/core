package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	"github.com/Team-Kujira/core/x/oracle/keeper"
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

// OracleQuery custom query interface for oracle querier
type OracleQuery struct {
	ExchangeRate *ExchangeRateQueryParams `json:"exchange_rate,omitempty"`
}

// ExchangeRateQueryResponse - exchange rates query response item
type ExchangeRateQueryResponse struct {
	Rate string `json:"rate"`
}

// QueryCustom implements custom query interface
func Handle(keeper keeper.Keeper, ctx sdk.Context, q *OracleQuery) (any, error) {

	if q.ExchangeRate != nil {
		rate, err := keeper.GetExchangeRate(ctx, q.ExchangeRate.Denom)
		if err != nil {
			return nil, err
		}

		return ExchangeRateQueryResponse{
			Rate: rate.String(),
		}, nil
	}

	return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown Oracle variant"}
}
