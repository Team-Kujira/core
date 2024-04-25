package cmd

// import (
// 	"errors"
// 	"io"
// 	"os"

// 	"cosmossdk.io/depinject"
// 	"cosmossdk.io/log"
// 	"cosmossdk.io/simapp"
// 	confixcmd "cosmossdk.io/tools/confix/cmd"

// 	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
// 	"github.com/Team-Kujira/core/app"
// 	cmtcfg "github.com/cometbft/cometbft/config"
// 	txmodule "github.com/cosmos/cosmos-sdk/x/auth/tx/config"
// 	"github.com/prometheus/client_golang/prometheus"

// 	// rosettaCmd "cosmossdk.io/tools/rosetta/cmd"
// 	"cosmossdk.io/client/v2/autocli"
// 	dbm "github.com/cosmos/cosmos-db"
// 	"github.com/cosmos/cosmos-sdk/client"
// 	"github.com/cosmos/cosmos-sdk/client/config"
// 	"github.com/cosmos/cosmos-sdk/client/debug"
// 	"github.com/cosmos/cosmos-sdk/client/flags"
// 	"github.com/cosmos/cosmos-sdk/client/keys"
// 	"github.com/cosmos/cosmos-sdk/client/pruning"
// 	"github.com/cosmos/cosmos-sdk/client/rpc"
// 	"github.com/cosmos/cosmos-sdk/client/snapshot"
// 	"github.com/cosmos/cosmos-sdk/codec"
// 	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
// 	"github.com/cosmos/cosmos-sdk/server"
// 	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
// 	servertypes "github.com/cosmos/cosmos-sdk/server/types"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/cosmos/cosmos-sdk/types/module"
// 	"github.com/cosmos/cosmos-sdk/types/tx/signing"
// 	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
// 	"github.com/cosmos/cosmos-sdk/x/auth/tx"
// 	"github.com/cosmos/cosmos-sdk/x/auth/types"
// 	"github.com/cosmos/cosmos-sdk/x/crisis"
// 	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
// 	"github.com/spf13/cast"
// 	"github.com/spf13/cobra"
// 	"github.com/spf13/viper"
// )

// // NewRootCmd creates a new root command for wasmd. It is called once in the
// // main function.
// func NewRootCmd() *cobra.Command {
// 	var (
// 		interfaceRegistry  codectypes.InterfaceRegistry
// 		appCodec           codec.Codec
// 		txConfig           client.TxConfig
// 		legacyAmino        *codec.LegacyAmino
// 		autoCliOpts        autocli.AppOptions
// 		moduleBasicManager module.BasicManager
// 	)

// 	if err := depinject.Inject(depinject.Configs(simapp.AppConfig, depinject.Supply(log.NewNopLogger())),
// 		&interfaceRegistry,
// 		&appCodec,
// 		&txConfig,
// 		&legacyAmino,
// 		&autoCliOpts,
// 		&moduleBasicManager,
// 	); err != nil {
// 		panic(err)
// 	}

// 	initClientCtx := client.Context{}.
// 		WithCodec(appCodec).
// 		WithInterfaceRegistry(interfaceRegistry).
// 		WithLegacyAmino(legacyAmino).
// 		WithInput(os.Stdin).
// 		WithAccountRetriever(types.AccountRetriever{}).
// 		WithHomeDir(simapp.DefaultNodeHome).
// 		WithViper("") // In simapp, we don't use any prefix for env variables.

// 	rootCmd := &cobra.Command{
// 		Use:   "simd",
// 		Short: "simulation app",
// 		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
// 			// set the default command outputs
// 			cmd.SetOut(cmd.OutOrStdout())
// 			cmd.SetErr(cmd.ErrOrStderr())

// 			initClientCtx = initClientCtx.WithCmdContext(cmd.Context())
// 			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
// 			if err != nil {
// 				return err
// 			}

// 			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
// 			if err != nil {
// 				return err
// 			}

// 			// This needs to go after ReadFromClientConfig, as that function
// 			// sets the RPC client needed for SIGN_MODE_TEXTUAL.
// 			enabledSignModes := append(tx.DefaultSignModes, signing.SignMode_SIGN_MODE_TEXTUAL)
// 			txConfigOpts := tx.ConfigOptions{
// 				EnabledSignModes:           enabledSignModes,
// 				TextualCoinMetadataQueryFn: txmodule.NewGRPCCoinMetadataQueryFn(initClientCtx),
// 			}
// 			txConfigWithTextual, err := tx.NewTxConfigWithOptions(
// 				codec.NewProtoCodec(interfaceRegistry),
// 				txConfigOpts,
// 			)
// 			if err != nil {
// 				return err
// 			}
// 			initClientCtx = initClientCtx.WithTxConfig(txConfigWithTextual)
// 			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
// 				return err
// 			}

// 			customAppTemplate, customAppConfig := initAppConfig()
// 			customCMTConfig := initCometBFTConfig()

// 			return server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, customCMTConfig)
// 		},
// 	}

// 	initRootCmd(rootCmd, txConfig, interfaceRegistry, appCodec, moduleBasicManager)

// 	if err := autoCliOpts.EnhanceRootCommand(rootCmd); err != nil {
// 		panic(err)
// 	}

// 	return rootCmd
// }

// // initCometBFTConfig helps to override default CometBFT Config values.
// // return cmtcfg.DefaultConfig if no custom configuration is required for the application.
// func initCometBFTConfig() *cmtcfg.Config {
// 	cfg := cmtcfg.DefaultConfig()

// 	// these values put a higher strain on node memory
// 	// cfg.P2P.MaxNumInboundPeers = 100
// 	// cfg.P2P.MaxNumOutboundPeers = 40

// 	return cfg
// }

// // initAppConfig helps to override default appConfig template and configs.
// // return "", nil if no custom configuration is required for the application.
// func initAppConfig() (string, interface{}) {
// 	// The following code snippet is just for reference.

// 	// WASMConfig defines configuration for the wasm module.
// 	type WASMConfig struct {
// 		// This is the maximum sdk gas (wasm and storage) that we allow for any x/wasm "smart" queries
// 		QueryGasLimit uint64 `mapstructure:"query_gas_limit"`

// 		// Address defines the gRPC-web server to listen on
// 		LruSize uint64 `mapstructure:"lru_size"`
// 	}

// 	type CustomAppConfig struct {
// 		serverconfig.Config

// 		WASM WASMConfig `mapstructure:"wasm"`
// 	}

// 	// Optionally allow the chain developer to overwrite the SDK's default
// 	// server config.
// 	srvCfg := serverconfig.DefaultConfig()
// 	// The SDK's default minimum gas price is set to "" (empty value) inside
// 	// app.toml. If left empty by validators, the node will halt on startup.
// 	// However, the chain developer can set a default app.toml value for their
// 	// validators here.
// 	//
// 	// In summary:
// 	// - if you leave srvCfg.MinGasPrices = "", all validators MUST tweak their
// 	//   own app.toml config,
// 	// - if you set srvCfg.MinGasPrices non-empty, validators CAN tweak their
// 	//   own app.toml to override, or use this default value.
// 	//
// 	// In simapp, we set the min gas prices to 0.
// 	srvCfg.MinGasPrices = "0stake"
// 	// srvCfg.BaseConfig.IAVLDisableFastNode = true // disable fastnode by default

// 	customAppConfig := CustomAppConfig{
// 		Config: *srvCfg,
// 		WASM: WASMConfig{
// 			LruSize:       1,
// 			QueryGasLimit: 300000,
// 		},
// 	}

// 	customAppTemplate := serverconfig.DefaultConfigTemplate + `
// [wasm]
// # This is the maximum sdk gas (wasm and storage) that we allow for any x/wasm "smart" queries
// query_gas_limit = 300000
// # This is the number of wasm vm instances we keep cached in memory for speed-up
// # Warning: this is currently unstable and may lead to crashes, best to keep for 0 unless testing locally
// lru_size = 0`

// 	return customAppTemplate, customAppConfig
// }

// func initRootCmd(
// 	rootCmd *cobra.Command,
// 	txConfig client.TxConfig,
// 	interfaceRegistry codectypes.InterfaceRegistry,
// 	appCodec codec.Codec,
// 	basicManager module.BasicManager,
// ) {
// 	cfg := sdk.GetConfig()
// 	cfg.Seal()

// 	rootCmd.AddCommand(
// 		genutilcli.InitCmd(basicManager, simapp.DefaultNodeHome),
// 		debug.Cmd(),
// 		confixcmd.ConfigCommand(),
// 		pruning.Cmd(newApp, simapp.DefaultNodeHome),
// 		snapshot.Cmd(newApp),
// 	)

// 	server.AddCommands(rootCmd, simapp.DefaultNodeHome, newApp, appExport, addModuleInitFlags)

// 	// add keybase, auxiliary RPC, query, genesis, and tx child commands
// 	rootCmd.AddCommand(
// 		server.StatusCommand(),
// 		genesisCommand(txConfig, basicManager),
// 		queryCommand(),
// 		txCommand(),
// 		keys.Commands(),
// 	)

// 	// add rosetta
// 	// rootCmd.AddCommand(rosettaCmd.RosettaCommand(interfaceRegistry, appCodec))
// }

// func addModuleInitFlags(startCmd *cobra.Command) {
// 	crisis.AddModuleInitFlags(startCmd)
// }

// func queryCommand() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:                        "query",
// 		Aliases:                    []string{"q"},
// 		Short:                      "Querying subcommands",
// 		DisableFlagParsing:         false,
// 		SuggestionsMinimumDistance: 2,
// 		RunE:                       client.ValidateCmd,
// 	}

// 	cmd.AddCommand(
// 		rpc.ValidatorCommand(),
// 		server.QueryBlockCmd(),
// 		authcmd.QueryTxsByEventsCmd(),
// 		server.QueryBlocksCmd(),
// 		authcmd.QueryTxCmd(),
// 	)

// 	return cmd
// }

// // genesisCommand builds genesis-related `simd genesis` command. Users may provide application specific commands as a parameter
// func genesisCommand(txConfig client.TxConfig, basicManager module.BasicManager, cmds ...*cobra.Command) *cobra.Command {
// 	cmd := genutilcli.Commands(txConfig, basicManager, simapp.DefaultNodeHome)

// 	for _, subCmd := range cmds {
// 		cmd.AddCommand(subCmd)
// 	}
// 	return cmd
// }

// func txCommand() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:                        "tx",
// 		Short:                      "Transactions subcommands",
// 		DisableFlagParsing:         false,
// 		SuggestionsMinimumDistance: 2,
// 		RunE:                       client.ValidateCmd,
// 	}

// 	cmd.AddCommand(
// 		authcmd.GetSignCommand(),
// 		authcmd.GetSignBatchCommand(),
// 		authcmd.GetMultiSignCommand(),
// 		authcmd.GetMultiSignBatchCmd(),
// 		authcmd.GetValidateSignaturesCommand(),
// 		authcmd.GetBroadcastCommand(),
// 		authcmd.GetEncodeCommand(),
// 		authcmd.GetDecodeCommand(),
// 		authcmd.GetSimulateCmd(),
// 	)

// 	return cmd
// }

// // newApp is an appCreator
// func newApp(
// 	logger log.Logger,
// 	db dbm.DB,
// 	traceStore io.Writer,
// 	appOpts servertypes.AppOptions,
// ) servertypes.Application {
// 	var wasmOpts []wasmkeeper.Option
// 	if cast.ToBool(appOpts.Get("telemetry.enabled")) {
// 		wasmOpts = append(wasmOpts, wasmkeeper.WithVMCacheMetrics(prometheus.DefaultRegisterer))
// 	}

// 	skipUpgradeHeights := make(map[int64]bool)
// 	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
// 		skipUpgradeHeights[int64(h)] = true
// 	}
// 	baseappOptions := server.DefaultBaseappOptions(appOpts)
// 	return app.New(
// 		logger,
// 		db,
// 		traceStore,
// 		true,
// 		appOpts,
// 		wasmOpts,
// 		baseappOptions...,
// 	)
// }

// // appExport creates a new kujiraApp (optionally at a given height)
// // and exports state.
// func appExport(
// 	logger log.Logger,
// 	db dbm.DB,
// 	traceStore io.Writer,
// 	height int64,
// 	forZeroHeight bool,
// 	jailAllowedAddrs []string,
// 	appOpts servertypes.AppOptions,
// 	modulesToExport []string,
// ) (servertypes.ExportedApp, error) {
// 	var kujiraApp *app.App

// 	// this check is necessary as we use the flag in x/upgrade.
// 	// we can exit more gracefully by checking the flag here.
// 	homePath, ok := appOpts.Get(flags.FlagHome).(string)
// 	if !ok || homePath == "" {
// 		return servertypes.ExportedApp{}, errors.New("application home not set")
// 	}

// 	viperAppOpts, ok := appOpts.(*viper.Viper)
// 	if !ok {
// 		return servertypes.ExportedApp{}, errors.New("appOpts is not viper.Viper")
// 	}

// 	// overwrite the FlagInvCheckPeriod
// 	viperAppOpts.Set(server.FlagInvCheckPeriod, 1)
// 	appOpts = viperAppOpts

// 	loadLatest := height == -1
// 	var emptyWasmOpts []wasmkeeper.Option
// 	kujiraApp = app.New(
// 		logger,
// 		db,
// 		traceStore,
// 		loadLatest,
// 		appOpts,
// 		emptyWasmOpts,
// 	)

// 	if height != -1 {
// 		if err := kujiraApp.LoadHeight(height); err != nil {
// 			return servertypes.ExportedApp{}, err
// 		}
// 	}

// 	return kujiraApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)
// }