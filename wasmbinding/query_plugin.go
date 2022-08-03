package wasmbinding

import (
	"encoding/json"
	"fmt"

	"kujira/wasmbinding/bindings"
	oracle "kujira/x/oracle/wasm"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// CustomQuerier dispatches custom CosmWasm bindings queries.
func CustomQuerier(qp *QueryPlugin) func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var contractQuery bindings.CosmosQuery
		if err := json.Unmarshal(request, &contractQuery); err != nil {
			return nil, sdkerrors.Wrap(err, "osmosis query")
		}

		if contractQuery.Oracle != nil {
			res, err := oracle.Handle(qp.oraclekeeper, ctx, contractQuery.Oracle)
			if err != nil {
				return nil, err
			}

			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
			}

			return bz, nil
		} else if contractQuery.Bank != nil {
			coin := qp.bankkeeper.GetSupply(ctx, contractQuery.Bank.Supply.Denom)
			res := banktypes.QuerySupplyOfResponse{
				Amount: coin,
			}

			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
			}

			return bz, nil
		} else if contractQuery.Denom != nil {
			var denomQuery = contractQuery.Denom
			switch {
			case denomQuery.FullDenom != nil:
				creator := denomQuery.FullDenom.CreatorAddr
				subdenom := denomQuery.FullDenom.Subdenom

				fullDenom, err := GetFullDenom(creator, subdenom)
				if err != nil {
					return nil, sdkerrors.Wrap(err, "osmo full denom query")
				}

				res := bindings.FullDenomResponse{
					Denom: fullDenom,
				}

				bz, err := json.Marshal(res)
				if err != nil {
					return nil, sdkerrors.Wrap(err, "osmo full denom query response")
				}

				return bz, nil

			case denomQuery.DenomAdmin != nil:
				res, err := qp.GetDenomAdmin(ctx, denomQuery.DenomAdmin.Subdenom)
				if err != nil {
					return nil, err
				}

				bz, err := json.Marshal(res)
				if err != nil {
					return nil, fmt.Errorf("failed to JSON marshal DenomAdminResponse response: %w", err)
				}

				return bz, nil

			default:
				return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown osmosis query variant"}
			}
		} else {
			return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown Custom variant"}
		}
	}
}

// ConvertSdkCoinsToWasmCoins converts sdk type coins to wasm vm type coins
func ConvertSdkCoinsToWasmCoins(coins []sdk.Coin) wasmvmtypes.Coins {
	var toSend wasmvmtypes.Coins
	for _, coin := range coins {
		c := ConvertSdkCoinToWasmCoin(coin)
		toSend = append(toSend, c)
	}
	return toSend
}

// ConvertSdkCoinToWasmCoin converts a sdk type coin to a wasm vm type coin
func ConvertSdkCoinToWasmCoin(coin sdk.Coin) wasmvmtypes.Coin {
	return wasmvmtypes.Coin{
		Denom: coin.Denom,
		// Note: gamm tokens have 18 decimal places, so 10^22 is common, no longer in u64 range
		Amount: coin.Amount.String(),
	}
}
