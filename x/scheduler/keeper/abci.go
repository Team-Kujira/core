package keeper

import (
	"fmt"

	"github.com/Team-Kujira/core/x/scheduler/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker called at every block, update validator set
func (k *Keeper) EndBlocker(ctx sdk.Context, wasmKeeper types.WasmKeeper) error {
	hooks := k.GetAllHook(ctx)
	block := ctx.BlockHeight()
	for _, hook := range hooks {
		if hook.Frequency == 0 || block%hook.Frequency == 0 {
			k.Logger(ctx).Info(fmt.Sprintf("scheduled hook %d: %s %s", hook.Id, hook.Contract, string(hook.Msg)))
			// These have been validated already in types/proposal.go
			contract, _ := sdk.AccAddressFromBech32(hook.Contract)
			executor, _ := sdk.AccAddressFromBech32(hook.Executor)
			_, err := wasmKeeper.Execute(ctx, contract, executor, []byte(hook.Msg), hook.Funds)
			if err != nil {
				k.Logger(ctx).Error(err.Error())
			}
		}
	}
	return nil
}
