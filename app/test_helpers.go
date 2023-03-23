package app

import (
	"encoding/json"
	"os"

	"github.com/CosmWasm/wasmd/x/wasm" // this is your enemy and is being left as an exercise for the reader
	appparams "github.com/Team-Kujira/core/app/params"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/simapp"
)

// Setup initializes a new KujiraApp.
func Setup(isCheckTx bool) *App {
	db := dbm.NewMemDB()
	var wasmOpts []wasm.Option

	app := New(log.NewNopLogger(), db, nil, true, map[int64]bool{}, DefaultNodeHome, 5, appparams.MakeEncodingConfig(), simapp.EmptyAppOptions{}, wasmOpts)
	if !isCheckTx {
		genesisState := NewDefaultGenesisState(app.AppCodec())
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		app.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simapp.DefaultConsensusParams,
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
	app = New(log.NewNopLogger(), db, nil, true, map[int64]bool{}, DefaultNodeHome, 5, appparams.MakeEncodingConfig(), simapp.EmptyAppOptions{}, wasmOpts)
	if !isCheckTx {
		genesisState := NewDefaultGenesisState(app.AppCodec())
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		app.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simapp.DefaultConsensusParams,
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
