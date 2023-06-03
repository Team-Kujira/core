package wasmbinding

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	"github.com/Team-Kujira/core/wasmbinding/bindings"
	denom "github.com/Team-Kujira/core/x/denom/wasm"
	oracle "github.com/Team-Kujira/core/x/oracle/wasm"

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
			return nil, errorsmod.Wrap(err, "kujira query")
		}

		switch {
		case contractQuery.Oracle != nil:
			res, err := oracle.Handle(qp.oraclekeeper, ctx, contractQuery.Oracle)
			if err != nil {
				return nil, err
			}

			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
			}

			return bz, nil

		case contractQuery.Bank != nil:
			if contractQuery.Bank.DenomMetadata != nil {
				metadata, _ := qp.bankkeeper.GetDenomMetaData(ctx, contractQuery.Bank.DenomMetadata.Denom)
				res := banktypes.QueryDenomMetadataResponse{
					Metadata: metadata,
				}

				bz, err := json.Marshal(res)
				if err != nil {
					return nil, errorsmod.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
				}

				return bz, nil
			}

			coin := qp.bankkeeper.GetSupply(ctx, contractQuery.Bank.Supply.Denom)
			res := banktypes.QuerySupplyOfResponse{
				Amount: coin,
			}

			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
			}

			return bz, nil

		case contractQuery.Denom != nil:
			res, err := denom.HandleQuery(qp.denomKeeper, ctx, contractQuery.Denom)
			if err != nil {
				return nil, err
			}

			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
			}

			return bz, nil

		default:
			return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown Custom variant"}
		}
	}
}
