package bindings

import (
	denom "github.com/Team-Kujira/core/x/denom/wasm"
	oracle "github.com/Team-Kujira/core/x/oracle/wasm"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// DenomQuery contains denom custom queries.

type CosmosQuery struct {
	Denom  *denom.DenomQuery
	Bank   *BankQuery
	Oracle *oracle.OracleQuery
}

type BankQuery struct {
	DenomMetadata *banktypes.QueryDenomMetadataRequest `json:"denom_metadata,omitempty"`
	Supply        *banktypes.QuerySupplyOfRequest      `json:"supply,omitempty"`
}
