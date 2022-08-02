package app

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	oraclekeeper "kujira/x/oracle/keeper"
	oracle "kujira/x/oracle/wasm"

	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// WasmQuerier - staking query interface for wasm contract
type WasmQuerier struct {
	bankkeeper   bankkeeper.Keeper
	oraclekeeper oraclekeeper.Keeper
}

// NewWasmQuerier return bank wasm query interface
func NewWasmQuerier(bankkeeper bankkeeper.Keeper, oraclekeeper oraclekeeper.Keeper) WasmQuerier {
	return WasmQuerier{bankkeeper, oraclekeeper}
}

// Query - implement query function
func (WasmQuerier) Query(_ sdk.Context, _ wasmvmtypes.QueryRequest) ([]byte, error) {
	return nil, nil
}

type BankQuery struct {
	Supply *banktypes.QuerySupplyOfRequest `json:"supply,omitempty"`
}

type CosmosQuery struct {
	Bank   *BankQuery
	Oracle *oracle.OracleQuery
}

// QueryCustom implements custom query interface
func (querier WasmQuerier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	var params CosmosQuery
	err := json.Unmarshal(data, &params)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var res any
	if params.Oracle != nil {
		res, err = oracle.Handle(querier.oraclekeeper, ctx, params.Oracle)
	} else if params.Bank != nil {
		coin := querier.bankkeeper.GetSupply(ctx, params.Bank.Supply.Denom)
		res = banktypes.QuerySupplyOfResponse{
			Amount: coin,
		}
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
