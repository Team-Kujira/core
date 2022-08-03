package bindings

import (
	denom "kujira/x/denom/wasm"
)

type CosmosMsg struct {
	Denom *denom.DenomMsg
}
