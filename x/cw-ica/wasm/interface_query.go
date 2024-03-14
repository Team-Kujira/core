package wasm

import (
	"cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Team-Kujira/core/x/cw-ica/keeper"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
)

// Querier - staking query interface for wasm contract
type Querier struct {
	keeper keeper.Keeper
}

// NewWasmQuerier return bank wasm query interface
func NewWasmQuerier(keeper keeper.Keeper) Querier {
	return Querier{keeper}
}

// Query - implement query function
func (Querier) Query(_ sdk.Context, _ wasmvmtypes.QueryRequest) ([]byte, error) {
	return nil, nil
}

type CwIcaQuery struct {
	/// Given the connection-id, owner, and account-id, returns the address
	/// of the interchain account.
	AccountAddress *AccountAddress `json:"account_address,omitempty"`
}

type AccountAddress struct {
	Owner        string `json:"owner"`
	ConnectionID string `json:"connection_id"`
	AccountID    string `json:"account_id"`
}

type AccountAddressResponse struct {
	Address string `json:"address"`
}

// QueryCustom implements custom query interface
func HandleQuery(k keeper.Keeper, ctx sdk.Context, q *CwIcaQuery) (any, error) {
	switch {
	case q.AccountAddress != nil:
		owner := q.AccountAddress.Owner + "-" + q.AccountAddress.AccountID

		portID, err := icatypes.NewControllerPortID(owner)
		if err != nil {
			return nil, errors.Wrap(err, "could not find account")
		}

		addr, found := k.IcaControllerKeeper().GetInterchainAccountAddress(ctx, q.AccountAddress.ConnectionID, portID)
		if !found {
			return nil, errors.Wrap(err, "no account found for portID")
		}

		return AccountAddressResponse{
			Address: addr,
		}, nil

	default:
		return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown CwIca Query variant"}
	}
}
