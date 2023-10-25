package keeper

import (
	"testing"

	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/Team-Kujira/core/x/scheduler/types"
	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
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
	ms := store.NewCommitMultiStore(db)
	key := sdk.NewKVStoreKey(types.StoreKey)
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

	ctx := sdk.NewContext(ms, tmproto.Header{}, false, log.NewNopLogger())
	keeper := NewKeeper(codec, key, authority)
	return keeper, ctx
}
