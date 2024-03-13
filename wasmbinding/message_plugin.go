package wasmbinding

import (
	"encoding/json"

	"cosmossdk.io/errors"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/terra-money/alliance/custom/bank/keeper"

	"github.com/Team-Kujira/core/wasmbinding/bindings"

	batchkeeper "github.com/Team-Kujira/core/x/batch/keeper"
	batch "github.com/Team-Kujira/core/x/batch/wasm"
	cwicakeeper "github.com/Team-Kujira/core/x/cw-ica/keeper"
	cwica "github.com/Team-Kujira/core/x/cw-ica/wasm"
	denomkeeper "github.com/Team-Kujira/core/x/denom/keeper"
	denom "github.com/Team-Kujira/core/x/denom/wasm"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
)

// CustomMessageDecorator returns decorator for custom CosmWasm bindings messages
func CustomMessageDecorator(
	bank bankkeeper.Keeper,
	denom denomkeeper.Keeper,
	cwica cwicakeeper.Keeper,
	ica icacontrollerkeeper.Keeper,
) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			wrapped: old,
			bank:    bank,
			denom:   denom,
			cwica:   cwica,
			ica:     ica,
		}
	}
}

type CustomMessenger struct {
	wrapped wasmkeeper.Messenger
	bank    bankkeeper.Keeper
	denom   denomkeeper.Keeper
	cwica   cwicakeeper.Keeper
	ica     icacontrollerkeeper.Keeper
	batch   batchkeeper.Keeper
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

		if contractMsg.Batch != nil {
			return batch.HandleMsg(m.batch, contractAddr, ctx, contractMsg.Batch)
		}

		if contractMsg.CwIca != nil {
			return cwica.HandleMsg(ctx, m.cwica, m.ica, contractAddr, contractMsg.CwIca)
		}

		return nil, nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown Custom variant"}
	}
	return m.wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
}
