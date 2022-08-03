package bindings

import (
	denom "kujira/x/denom/wasm"
	oracle "kujira/x/oracle/wasm"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// DenomQuery contains denom custom queries.

type CosmosQuery struct {
	Denom  *denom.DenomQuery
	Bank   *BankQuery
	Oracle *oracle.OracleQuery
}

type BankQuery struct {
	Supply *banktypes.QuerySupplyOfRequest `json:"supply,omitempty"`
}
