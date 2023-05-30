package bindings

import (
	denom "github.com/Team-Kujira/core/x/denom/wasm"
	intertx "github.com/Team-Kujira/core/x/inter-tx/wasm"
	oracle "github.com/Team-Kujira/core/x/oracle/wasm"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// DenomQuery contains denom custom queries.

type CosmosQuery struct {
	Denom   *denom.DenomQuery
	Bank    *BankQuery
	Oracle  *oracle.OracleQuery
	Intertx *intertx.IntertxQuery
}

type BankQuery struct {
	DenomMetadata *banktypes.QueryDenomMetadataRequest `json:"denom_metadata,omitempty"`
	Supply        *banktypes.QuerySupplyOfRequest      `json:"supply,omitempty"`
}
