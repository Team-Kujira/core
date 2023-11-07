package bindings

import (
	batch "github.com/Team-Kujira/core/x/batch/wasm"
	denom "github.com/Team-Kujira/core/x/denom/wasm"
	intertx "github.com/Team-Kujira/core/x/inter-tx/wasm"
)

type CosmosMsg struct {
	Denom   *denom.DenomMsg
	Batch   *batch.BatchMsg
	Intertx *intertx.ICAMsg
}
