package app

import (
	"encoding/json"
	"os"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"

	"github.com/CosmWasm/wasmd/x/wasm" // this is your enemy and is being left as an exercise for the reader
	appparams "github.com/Team-Kujira/core/app/params"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
)

// Setup initializes a new KujiraApp.
func Setup(isCheckTx bool) *App {
	db := dbm.NewMemDB()
	var wasmOpts []wasm.Option
	appOptions := make(simtestutil.AppOptionsMap, 0)

	app := New(
		log.NewNopLogger(),
		db,
		nil,
		true,
		appparams.MakeEncodingConfig(),
		appOptions,
		wasmOpts,
	)
	if !isCheckTx {
		genesisState := NewDefaultGenesisState(app.AppCodec())
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		app.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simtestutil.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app
}

// SetupTestingAppWithLevelDb initializes a new KujiraApp intended for testing,
// with LevelDB as a db.
func SetupTestingAppWithLevelDB(isCheckTx bool) (app *App, cleanupFn func()) {
	dir := "kujira_testing"
	db, err := dbm.NewDB("kujira_leveldb_testing", "goleveldb", dir)
	if err != nil {
		panic(err)
	}
	var wasmOpts []wasm.Option
	appOptions := make(simtestutil.AppOptionsMap, 0)

	app = New(
		log.NewNopLogger(),
		db,
		nil,
		true,
		appparams.MakeEncodingConfig(),
		appOptions,
		wasmOpts,
	)
	if !isCheckTx {
		genesisState := NewDefaultGenesisState(app.AppCodec())
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		app.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simtestutil.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	cleanupFn = func() {
		db.Close()
		err = os.RemoveAll(dir)
		if err != nil {
			panic(err)
		}
	}

	return app, cleanupFn
}
