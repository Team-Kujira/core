package app

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	oraclekeeper "kujira/x/oracle/keeper"
	oracle "kujira/x/oracle/wasm"
)

// WasmQuerier - staking query interface for wasm contract
type WasmQuerier struct {
	oraclekeeper oraclekeeper.Keeper
}

// NewWasmQuerier return bank wasm query interface
func NewWasmQuerier(oraclekeeper oraclekeeper.Keeper) WasmQuerier {
	return WasmQuerier{oraclekeeper}
}

// Query - implement query function
func (WasmQuerier) Query(_ sdk.Context, _ wasmvmtypes.QueryRequest) ([]byte, error) {
	return nil, nil
}

type CosmosQuery struct {
	Oracle *oracle.OracleQuery
}

// QueryCustom implements custom query interface
func (querier WasmQuerier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	var params CosmosQuery
	err := json.Unmarshal(data, &params)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var (
		res any
	)

	if params.Oracle != nil {
		res, err = oracle.Handle(querier.oraclekeeper, ctx, params.Oracle)
	} else {
		return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown Custom variant"}
	}

	if err != nil {
		return nil, err
	}

	bz, err := json.Marshal(res)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
