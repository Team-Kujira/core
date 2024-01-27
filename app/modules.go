package app

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	appparams "github.com/Team-Kujira/core/app/params"
	"github.com/Team-Kujira/core/x/denom"
	denomtypes "github.com/Team-Kujira/core/x/denom/types"
	"github.com/Team-Kujira/core/x/oracle"
	oracletypes "github.com/Team-Kujira/core/x/oracle/types"
	scheduler "github.com/Team-Kujira/core/x/scheduler"
	schedulertypes "github.com/Team-Kujira/core/x/scheduler/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	ibcfee "github.com/cosmos/ibc-go/v7/modules/apps/29-fee"
	ibcfeetypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	transfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	bank "github.com/terra-money/alliance/custom/bank"
	alliancemodule "github.com/terra-money/alliance/x/alliance"
	alliancemoduletypes "github.com/terra-money/alliance/x/alliance/types"
)

// module account permissions
var maccPerms = map[string][]string{
	authtypes.FeeCollectorName:          nil,
	distrtypes.ModuleName:               nil,
	minttypes.ModuleName:                {authtypes.Minter},
	stakingtypes.BondedPoolName:         {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName:      {authtypes.Burner, authtypes.Staking},
	govtypes.ModuleName:                 {authtypes.Burner},
	ibctransfertypes.ModuleName:         {authtypes.Minter, authtypes.Burner},
	ibcfeetypes.ModuleName:              nil,
	icatypes.ModuleName:                 nil,
	wasmtypes.ModuleName:                {authtypes.Burner},
	denomtypes.ModuleName:               {authtypes.Minter, authtypes.Burner},
	schedulertypes.ModuleName:           nil,
	oracletypes.ModuleName:              nil,
	alliancemoduletypes.ModuleName:      {authtypes.Minter, authtypes.Burner},
	alliancemoduletypes.RewardsPoolName: nil,
}

// ModuleBasics defines the module BasicManager is in charge of setting up basic,
// non-dependant module elements, such as codec registration
// and genesis verification.
var ModuleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
	authzmodule.AppModuleBasic{},
	bank.AppModule{},
	capability.AppModuleBasic{},
	staking.AppModuleBasic{},
	mint.AppModuleBasic{},
	distr.AppModuleBasic{},
	gov.NewAppModuleBasic(getGovProposalHandlers()),
	params.AppModuleBasic{},
	consensus.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	feegrantmodule.AppModuleBasic{},
	ibc.AppModuleBasic{},
	ibctm.AppModuleBasic{},

	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},
	transfer.AppModuleBasic{},
	vesting.AppModuleBasic{},
	wasm.AppModuleBasic{},
	ica.AppModuleBasic{},
	ibcfee.AppModuleBasic{},
	denom.AppModuleBasic{},
	scheduler.AppModuleBasic{},
	oracle.AppModuleBasic{},
	alliancemodule.AppModuleBasic{},
)

func appModules(
	app *App,
	encodingConfig appparams.EncodingConfig,
	skipGenesisInvariants bool,
) []module.AppModule {
	appCodec := encodingConfig.Codec

	return []module.AppModule{
		genutil.NewAppModule(
			app.AccountKeeper,
			app.StakingKeeper,
			app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),

		auth.NewAppModule(
			appCodec,
			app.AccountKeeper,
			authsims.RandomGenesisAccounts,
			app.GetSubspace(authtypes.ModuleName),
		),

		authzmodule.NewAppModule(
			appCodec,
			app.AuthzKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.interfaceRegistry,
		),

		vesting.NewAppModule(
			app.AccountKeeper,
			app.BankKeeper,
		),

		bank.NewAppModule(
			appCodec,
			app.BankKeeper,
			app.AccountKeeper,
			app.GetSubspace(banktypes.ModuleName),
		),

		capability.NewAppModule(
			appCodec,
			*app.CapabilityKeeper,
			false,
		),

		feegrantmodule.NewAppModule(
			appCodec,
			app.AccountKeeper,
			app.BankKeeper,
			app.FeeGrantKeeper,
			app.interfaceRegistry,
		),

		gov.NewAppModule(
			appCodec,
			&app.GovKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.GetSubspace(govtypes.ModuleName),
		),

		mint.NewAppModule(
			appCodec,
			app.MintKeeper,
			app.AccountKeeper,
			nil,
			app.GetSubspace(minttypes.ModuleName),
		),

		slashing.NewAppModule(
			appCodec,
			app.SlashingKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.StakingKeeper,
			app.GetSubspace(slashingtypes.ModuleName),
		),

		distr.NewAppModule(
			appCodec,
			app.DistrKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.StakingKeeper,
			app.GetSubspace(distrtypes.ModuleName),
		),

		staking.NewAppModule(
			appCodec,
			app.StakingKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.GetSubspace(stakingtypes.ModuleName),
		),

		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		params.NewAppModule(app.ParamsKeeper),
		transfer.NewAppModule(app.TransferKeeper),

		ibcfee.NewAppModule(app.IBCFeeKeeper),

		consensus.NewAppModule(appCodec, app.ConsensusParamsKeeper),

		wasm.NewAppModule(
			appCodec,
			&app.WasmKeeper,
			app.StakingKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.MsgServiceRouter(),
			app.GetSubspace(wasmtypes.ModuleName),
		),

		denom.NewAppModule(
			appCodec,
			*app.DenomKeeper,
			app.AccountKeeper,
			app.BankKeeper,
		),

		ica.NewAppModule(&app.ICAControllerKeeper, &app.ICAHostKeeper),

		scheduler.NewAppModule(
			appCodec,
			app.SchedulerKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			wasmkeeper.NewDefaultPermissionKeeper(app.WasmKeeper),
		),

		oracle.NewAppModule(
			appCodec,
			app.OracleKeeper,
			app.AccountKeeper,
			app.BankKeeper,
		),

		alliancemodule.NewAppModule(
			appCodec,
			app.AllianceKeeper,
			app.StakingKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.interfaceRegistry,
			app.GetSubspace(alliancemoduletypes.ModuleName),
		),

		crisis.NewAppModule(
			app.CrisisKeeper,
			skipGenesisInvariants,
			app.GetSubspace(crisistypes.ModuleName),
		),
	}
}

// orderBeginBlockers tell the app's module manager how to set the order of
// BeginBlockers, which are run at the beginning of every block.
func orderBeginBlockers() []string {
	return []string{
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		minttypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		vestingtypes.ModuleName,
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		authtypes.ModuleName,
		authz.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		consensusparamtypes.ModuleName,
		icatypes.ModuleName,
		ibcfeetypes.ModuleName,

		wasmtypes.ModuleName,
		denomtypes.ModuleName,
		schedulertypes.ModuleName,
		oracletypes.ModuleName,
		alliancemoduletypes.ModuleName,
	}
}

func orderEndBlockers() []string {
	return []string{
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		authz.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		vestingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		consensusparamtypes.ModuleName,
		icatypes.ModuleName,
		ibcfeetypes.ModuleName,

		wasmtypes.ModuleName,
		denomtypes.ModuleName,
		schedulertypes.ModuleName,
		oracletypes.ModuleName,
		alliancemoduletypes.ModuleName,
	}
}

func orderInitBlockers() []string {
	return []string{
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		authz.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		vestingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		ibcexported.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		feegrant.ModuleName,
		ibctransfertypes.ModuleName,
		consensusparamtypes.ModuleName,
		icatypes.ModuleName,
		ibcfeetypes.ModuleName,

		denomtypes.ModuleName,
		schedulertypes.ModuleName,
		oracletypes.ModuleName,
		alliancemoduletypes.ModuleName,
		wasmtypes.ModuleName,
	}
}
