package bindings

import (
	denom "github.com/Team-Kujira/core/x/denom/wasm"
	intertx "github.com/Team-Kujira/core/x/inter-tx/wasm"
)

type CosmosMsg struct {
	Denom   *denom.DenomMsg
	Intertx *intertx.ICAMsg
}
