package bindings

import (
	oracle "kujira/x/oracle/wasm"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// DenomQuery contains denom custom queries.

type CosmosQuery struct {
	Denom  *DenomQuery
	Bank   *BankQuery
	Oracle *oracle.OracleQuery
}

type DenomQuery struct {
	/// Given a subdenom minted by a contract via `DenomMsg::MintTokens`,
	/// returns the full denom as used by `BankMsg::Send`.
	FullDenom *FullDenom `json:"full_denom,omitempty"`
	/// Returns the admin of a denom, if the denom is a Token Factory denom.
	DenomAdmin *DenomAdmin `json:"denom_admin,omitempty"`
}

type BankQuery struct {
	Supply *banktypes.QuerySupplyOfRequest `json:"supply,omitempty"`
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
