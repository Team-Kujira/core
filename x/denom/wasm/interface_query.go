package wasm

import (
	"fmt"

	"cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Team-Kujira/core/x/denom/keeper"
	denomtypes "github.com/Team-Kujira/core/x/denom/types"
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

type DenomQuery struct {
	/// Given a subdenom minted by a contract via `DenomMsg::MintTokens`,
	/// returns the full denom as used by `BankMsg::Send`.
	FullDenom *FullDenom `json:"full_denom,omitempty"`
	/// Returns the admin of a denom, if the denom is a Token Factory denom.
	DenomAdmin *DenomAdmin `json:"denom_admin,omitempty"`
}

type FullDenom struct {
	CreatorAddr string `json:"creator_addr"`
	Subdenom    string `json:"subdenom"`
}

type DenomAdmin struct {
	Subdenom string `json:"subdenom"`
}

type FullDenomResponse struct {
	Denom string `json:"denom"`
}

type DenomAdminResponse struct {
	Admin string `json:"admin"`
}

// GetFullDenom is a function, not method, so the message_plugin can use it
func GetFullDenom(contract string, subDenom string) (string, error) {
	// Address validation
	if _, err := parseAddress(contract); err != nil {
		return "", err
	}
	fullDenom, err := denomtypes.GetTokenDenom(contract, subDenom)
	if err != nil {
		return "", errors.Wrap(err, "validate sub-denom")
	}

	return fullDenom, nil
}

// parseAddress parses address from bech32 string and verifies its format.
func parseAddress(addr string) (sdk.AccAddress, error) {
	parsed, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return nil, errors.Wrap(err, "address from bech32")
	}
	err = sdk.VerifyAddressFormat(parsed)
	if err != nil {
		return nil, errors.Wrap(err, "verify address format")
	}
	return parsed, nil
}

// GetDenomAdmin is a query to get denom admin.
func GetDenomAdmin(keeper keeper.Keeper, ctx sdk.Context, denom string) (*DenomAdminResponse, error) {
	metadata, err := keeper.GetAuthorityMetadata(ctx, denom)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin for denom: %s", denom)
	}

	return &DenomAdminResponse{Admin: metadata.Admin}, nil
}

// QueryCustom implements custom query interface
func HandleQuery(keeper keeper.Keeper, ctx sdk.Context, q *DenomQuery) (any, error) {
	switch {
	case q.FullDenom != nil:
		creator := q.FullDenom.CreatorAddr
		subdenom := q.FullDenom.Subdenom

		fullDenom, err := GetFullDenom(creator, subdenom)
		if err != nil {
			return nil, err
		}

		return FullDenomResponse{
			Denom: fullDenom,
		}, nil

	case q.DenomAdmin != nil:
		res, err := GetDenomAdmin(keeper, ctx, q.DenomAdmin.Subdenom)
		if err != nil {
			return nil, err
		}

		return res, nil

	default:
		return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown Denom variant"}
	}
}
