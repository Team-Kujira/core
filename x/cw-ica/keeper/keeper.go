package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/Team-Kujira/core/x/cw-ica/types"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/keeper"
)

type Keeper struct {
	Codec codec.Codec

	storeKey storetypes.StoreKey

	scopedKeeper        capabilitykeeper.ScopedKeeper
	icaControllerKeeper icacontrollerkeeper.Keeper
	wasmKeeper          *wasmkeeper.Keeper
}

func NewKeeper(cdc codec.Codec, storeKey storetypes.StoreKey, iaKeeper icacontrollerkeeper.Keeper, scopedKeeper capabilitykeeper.ScopedKeeper, wasmKeeper *wasmkeeper.Keeper) Keeper {
	return Keeper{
		Codec:    cdc,
		storeKey: storeKey,

		scopedKeeper:        scopedKeeper,
		icaControllerKeeper: iaKeeper,
		wasmKeeper:          wasmKeeper,
	}
}

// Logger returns the application logger, scoped to the associated module
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// ClaimCapability claims the channel capability passed via the OnOpenChanInit callback
func (k *Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}

func (k Keeper) IcaControllerKeeper() icacontrollerkeeper.Keeper {
	return k.icaControllerKeeper
}
