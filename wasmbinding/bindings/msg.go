package bindings

import (
	batch "github.com/Team-Kujira/core/x/batch/wasm"
	cwica "github.com/Team-Kujira/core/x/cw-ica/wasm"
	denom "github.com/Team-Kujira/core/x/denom/wasm"
)

type CosmosMsg struct {
	Denom *denom.DenomMsg
	Batch *batch.BatchMsg
	CwIca *cwica.ICAMsg
}
