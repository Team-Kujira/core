package wasmbinding

import (
	"encoding/json"

	"cosmossdk.io/errors"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/terra-money/alliance/custom/bank/keeper"

	"github.com/Team-Kujira/core/wasmbinding/bindings"

	denom "github.com/Team-Kujira/core/x/denom/wasm"

	denomkeeper "github.com/Team-Kujira/core/x/denom/keeper"
)

// CustomMessageDecorator returns decorator for custom CosmWasm bindings messages
func CustomMessageDecorator(
	bank bankkeeper.Keeper,
	denom denomkeeper.Keeper,
) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			wrapped: old,
			bank:    bank,
			denom:   denom,
		}
	}
}

type CustomMessenger struct {
	wrapped wasmkeeper.Messenger
	bank    bankkeeper.Keeper
	denom   denomkeeper.Keeper
}

var _ wasmkeeper.Messenger = (*CustomMessenger)(nil)

// DispatchMsg executes on the contractMsg.
func (m *CustomMessenger) DispatchMsg(
	ctx sdk.Context,
	contractAddr sdk.AccAddress,
	contractIBCPortID string,
	msg wasmvmtypes.CosmosMsg,
) ([]sdk.Event, [][]byte, error) {
	if msg.Custom != nil {
		// only handle the happy path where this is really creating / minting / swapping ...
		// leave everything else for the wrapped version
		var contractMsg bindings.CosmosMsg
		if err := json.Unmarshal(msg.Custom, &contractMsg); err != nil {
			return nil, nil, errors.Wrap(err, "kujira msg")
		}

		if contractMsg.Denom != nil {
			return denom.HandleMsg(m.denom, m.bank, contractAddr, ctx, contractMsg.Denom)
		}

		return nil, nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown Custom variant"}
	}
	return m.wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
}
