package bindings

import (
	denom "github.com/Team-Kujira/core/x/denom/wasm"
)

type CosmosMsg struct {
	Denom *denom.DenomMsg
}
