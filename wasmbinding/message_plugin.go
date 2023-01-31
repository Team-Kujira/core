package wasmbinding

import (
	"encoding/json"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	bankkeeper "github.com/terra-money/alliance/custom/bank/keeper"

	"github.com/Team-Kujira/core/wasmbinding/bindings"

	denom "github.com/Team-Kujira/core/x/denom/wasm"

	denomkeeper "github.com/Team-Kujira/core/x/denom/keeper"
)

// CustomMessageDecorator returns decorator for custom CosmWasm bindings messages
func CustomMessageDecorator(
	bank bankkeeper.Keeper,
	denom denomkeeper.Keeper,
	auth authkeeper.AccountKeeper,
) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			wrapped: old,
			auth:    auth,
			bank:    bank,
			denom:   denom,
		}
	}
}

type CustomMessenger struct {
	wrapped wasmkeeper.Messenger
	auth    authkeeper.AccountKeeper
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
			return nil, nil, sdkerrors.Wrap(err, "kujira msg")
		}

		if contractMsg.Denom != nil {
			return denom.HandleMsg(m.denom, m.bank, contractAddr, ctx, contractMsg.Denom)
		} else if contractMsg.Auth != nil {
			handler := bindings.AuthHandler{AccountKeeper: m.auth, BankKeeper: m.bank}
			cva := contractMsg.Auth.CreateVestingAccount
			return handler.CreateVestingAccount(ctx, &types.MsgCreateVestingAccount{
				FromAddress: contractAddr.String(),
				ToAddress:   cva.ToAddress,
				Amount:      cva.Amount,
				// Timestamp is picoseconds. MsgCreateVestingAccount expects seconds
				EndTime: int64(cva.EndTime.Uint64() / 1000000000),
				Delayed: cva.Delayed,
			})
		} else {
			return nil, nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown Custom variant"}
		}
	}
	return m.wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
}
