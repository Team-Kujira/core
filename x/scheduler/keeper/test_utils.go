package keeper

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/Team-Kujira/core/x/scheduler/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"
)

var ModuleBasics = module.NewBasicManager(
	bank.AppModuleBasic{},
	wasm.AppModuleBasic{},
)

// Setup the testing environment
func CreateTestKeeper(t *testing.T) (Keeper, sdk.Context) {
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)

	ms := store.NewCommitMultiStore(db, logger, storemetrics.NewNoOpMetrics())
	key := storetypes.NewKVStoreKey(types.StoreKey)
	ms.MountStoreWithDB(key, storetypes.StoreTypeIAVL, db)
	require.NoError(t, ms.LoadLatestVersion())

	amino := codec.NewLegacyAmino()
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	codec := codec.NewProtoCodec(interfaceRegistry)
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	std.RegisterInterfaces(interfaceRegistry)
	std.RegisterLegacyAminoCodec(amino)

	ModuleBasics.RegisterLegacyAminoCodec(amino)
	ModuleBasics.RegisterInterfaces(interfaceRegistry)
	types.RegisterLegacyAminoCodec(amino)
	types.RegisterInterfaces(interfaceRegistry)

	ctx := sdk.NewContext(ms, tmproto.Header{}, false, nil)
	keeper := NewKeeper(codec, key, authority)
	return keeper, ctx
}
