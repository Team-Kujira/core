package bindings

import (
	batch "github.com/Team-Kujira/core/x/batch/wasm"
	denom "github.com/Team-Kujira/core/x/denom/wasm"
)

type CosmosMsg struct {
	Denom *denom.DenomMsg
	Batch *batch.BatchMsg
}
