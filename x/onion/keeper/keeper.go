package keeper

import (
	"fmt"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/Team-Kujira/core/x/onion/types"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/baseapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type (
	Keeper struct {
		storeKey   storetypes.StoreKey
		paramSpace paramtypes.Subspace

		accountKeeper types.AccountKeeper

		router          *baseapp.MsgServiceRouter
		signModeHandler authsigning.SignModeHandler
	}
)

// NewKeeper returns a new instance of the x/ibchooks keeper
func NewKeeper(
	storeKey storetypes.StoreKey,
	paramSpace paramtypes.Subspace,
	contractKeeper *wasmkeeper.PermissionedKeeper,
	router *baseapp.MsgServiceRouter,
	accountKeeper types.AccountKeeper,
	signModeHandler authsigning.SignModeHandler,
) *Keeper {
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}
	return &Keeper{
		storeKey:        storeKey,
		paramSpace:      paramSpace,
		router:          router,
		accountKeeper:   accountKeeper,
		signModeHandler: signModeHandler,
	}
}

// Logger returns a logger for the x/tokenfactory module
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetParams returns the total set of the module's parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the module's parameters with the provided parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	k.SetParams(ctx, genState.Params)
	for _, seq := range genState.Sequences {
		k.SetSequence(ctx, seq)
	}
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params:    k.GetParams(ctx),
		Sequences: k.GetAllSequences(ctx),
	}
}
