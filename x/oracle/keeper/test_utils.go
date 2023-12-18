//nolint:all
package keeper

import (
	"testing"
	"time"

	"github.com/Team-Kujira/core/x/oracle/types"

	storemetrics "cosmossdk.io/store/metrics"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/std"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	auth "github.com/cosmos/cosmos-sdk/x/auth"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	bank "github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	params "github.com/cosmos/cosmos-sdk/x/params"
	staking "github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/log"
	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/secp256k1"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"

	"cosmossdk.io/math"
	"cosmossdk.io/store"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

const faucetAccountName = "faucet"

// ModuleBasics nolint
var ModuleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	bank.AppModuleBasic{},
	distr.AppModuleBasic{},
	staking.AppModuleBasic{},
	params.AppModuleBasic{},
)

// Test addresses
var (
	ValPubKeys = simtestutil.CreateTestPubKeys(5)

	pubKeys = []crypto.PubKey{
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
	}

	Addrs = []sdk.AccAddress{
		sdk.AccAddress(pubKeys[0].Address()),
		sdk.AccAddress(pubKeys[1].Address()),
		sdk.AccAddress(pubKeys[2].Address()),
		sdk.AccAddress(pubKeys[3].Address()),
		sdk.AccAddress(pubKeys[4].Address()),
	}

	ValAddrs = []sdk.ValAddress{
		sdk.ValAddress(pubKeys[0].Address()),
		sdk.ValAddress(pubKeys[1].Address()),
		sdk.ValAddress(pubKeys[2].Address()),
		sdk.ValAddress(pubKeys[3].Address()),
		sdk.ValAddress(pubKeys[4].Address()),
	}

	InitTokens = sdk.TokensFromConsensusPower(200, sdk.DefaultPowerReduction)
	InitCoins  = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, InitTokens))

	OracleDecPrecision = 8

	Bech32Prefix         = "cosmos"
	Bech32PrefixValAddr  = Bech32Prefix + "valoper"
	Bech32PrefixConsAddr = Bech32Prefix + "valcons"
)

// TestInput nolint
type TestInput struct {
	Ctx           sdk.Context
	Cdc           *codec.LegacyAmino
	AccountKeeper authkeeper.AccountKeeper
	BankKeeper    bankkeeper.Keeper
	OracleKeeper  Keeper
	StakingKeeper stakingkeeper.Keeper
	DistrKeeper   distrkeeper.Keeper
}

// CreateTestInput nolint
func CreateTestInput(t *testing.T) TestInput {
	keyAcc := storetypes.NewKVStoreKey(authtypes.StoreKey)
	keyBank := storetypes.NewKVStoreKey(banktypes.StoreKey)
	keyParams := storetypes.NewKVStoreKey(paramstypes.StoreKey)
	tKeyParams := storetypes.NewTransientStoreKey(paramstypes.TStoreKey)
	keyOracle := storetypes.NewKVStoreKey(types.StoreKey)
	keySlashing := storetypes.NewKVStoreKey(slashingtypes.StoreKey)
	keyStaking := storetypes.NewKVStoreKey(stakingtypes.StoreKey)
	keyDistr := storetypes.NewKVStoreKey(distrtypes.StoreKey)
	logger := log.NewTestLogger(t)
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db, logger, storemetrics.NewNoOpMetrics())
	ctx := sdk.NewContext(ms, tmproto.Header{Time: time.Now().UTC()}, false, log.NewNopLogger())

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)

	appCodec := codec.NewProtoCodec(interfaceRegistry)
	legacyAmino := codec.NewLegacyAmino()

	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	ms.MountStoreWithDB(keyAcc, storetypes.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyBank, storetypes.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyParams, storetypes.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyParams, storetypes.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyOracle, storetypes.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySlashing, storetypes.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, storetypes.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyDistr, storetypes.StoreTypeIAVL, db)

	require.NoError(t, ms.LoadLatestVersion())

	blackListAddrs := map[string]bool{
		authtypes.FeeCollectorName:     true,
		stakingtypes.NotBondedPoolName: true,
		stakingtypes.BondedPoolName:    true,
		distrtypes.ModuleName:          true,
		faucetAccountName:              true,
	}

	maccPerms := map[string][]string{
		faucetAccountName:              {authtypes.Minter},
		authtypes.FeeCollectorName:     nil,
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		distrtypes.ModuleName:          nil,
		types.ModuleName:               nil,
	}

	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, keyParams, tKeyParams)
	accountKeeper := authkeeper.NewAccountKeeper(
		appCodec,
		runtime.NewKVStoreService(keyAcc),
		authtypes.ProtoBaseAccount,
		maccPerms,
		authcodec.NewBech32Codec(Bech32Prefix),
		Bech32Prefix,
		authority,
	)
	bankKeeper := bankkeeper.NewBaseKeeper(
		appCodec,
		runtime.NewKVStoreService(keyBank),
		accountKeeper,
		blackListAddrs,
		authority,
		logger,
	)

	totalSupply := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, InitTokens.MulRaw(int64(len(Addrs)*10))))
	bankKeeper.MintCoins(ctx, faucetAccountName, totalSupply)

	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keyStaking),
		accountKeeper,
		bankKeeper,
		authority,
		authcodec.NewBech32Codec(Bech32PrefixValAddr),
		authcodec.NewBech32Codec(Bech32PrefixConsAddr),
	)

	stakingParams := stakingtypes.DefaultParams()
	stakingKeeper.SetParams(ctx, stakingParams)

	slashingKeeper := slashingkeeper.NewKeeper(
		appCodec,
		legacyAmino,
		runtime.NewKVStoreService(keySlashing),
		stakingKeeper,
		authority,
	)

	distrKeeper := distrkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keyDistr),
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		authtypes.FeeCollectorName,
		authority,
	)

	distrKeeper.FeePool.Set(ctx, distrtypes.InitialFeePool())
	distrParams := distrtypes.DefaultParams()
	distrParams.CommunityTax = math.LegacyNewDecWithPrec(2, 2)
	distrParams.BaseProposerReward = math.LegacyNewDecWithPrec(1, 2)
	distrParams.BonusProposerReward = math.LegacyNewDecWithPrec(4, 2)
	distrKeeper.Params.Set(ctx, distrParams)
	stakingKeeper.SetHooks(stakingtypes.NewMultiStakingHooks(distrKeeper.Hooks()))

	bankKeeper.SendCoinsFromModuleToModule(ctx, faucetAccountName, stakingtypes.NotBondedPoolName, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, InitTokens.MulRaw(int64(len(Addrs))))))

	for _, addr := range Addrs {
		err := bankKeeper.SendCoinsFromModuleToAccount(ctx, faucetAccountName, addr, InitCoins)
		require.NoError(t, err)
	}

	keeper := NewKeeper(
		appCodec,
		keyOracle,
		nil,
		paramsKeeper.Subspace(types.ModuleName),
		accountKeeper,
		bankKeeper,
		distrKeeper,
		slashingKeeper,
		stakingKeeper,
		distrtypes.ModuleName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	defaults := types.DefaultParams()
	keeper.SetParams(ctx, defaults)

	return TestInput{ctx, legacyAmino, accountKeeper, bankKeeper, keeper, *stakingKeeper, distrKeeper}
}

// NewTestMsgCreateValidator test msg creator
func NewTestMsgCreateValidator(address sdk.ValAddress, pubKey cryptotypes.PubKey, amt math.Int) *stakingtypes.MsgCreateValidator {
	commission := stakingtypes.NewCommissionRates(math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec())
	msg, _ := stakingtypes.NewMsgCreateValidator(
		address.String(), pubKey, sdk.NewCoin(sdk.DefaultBondDenom, amt),
		stakingtypes.Description{Moniker: "moniker"}, commission, math.OneInt(),
	)

	return msg
}

// FundAccount is a utility function that funds an account by minting and
// sending the coins to the address. This should be used for testing purposes
// only!
func FundAccount(input TestInput, addr sdk.AccAddress, amounts sdk.Coins) error {
	if err := input.BankKeeper.MintCoins(input.Ctx, faucetAccountName, amounts); err != nil {
		return err
	}

	return input.BankKeeper.SendCoinsFromModuleToAccount(input.Ctx, faucetAccountName, addr, amounts)
}
