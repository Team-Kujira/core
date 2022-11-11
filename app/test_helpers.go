package app

import (
	"encoding/json"
	"os"

	"github.com/ignite-hq/cli/ignite/pkg/cosmoscmd"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Setup initializes a new KujiraApp.
func Setup(isCheckTx bool) *App {
	db := dbm.NewMemDB()
	app := New(log.NewNopLogger(), db, nil, true, map[int64]bool{}, DefaultNodeHome, 5, cosmoscmd.MakeEncodingConfig(ModuleBasics), simapp.EmptyAppOptions{})
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
	db, err := sdk.NewLevelDB("kujira_leveldb_testing", dir)
	if err != nil {
		panic(err)
	}
	app = New(log.NewNopLogger(), db, nil, true, map[int64]bool{}, DefaultNodeHome, 5, cosmoscmd.MakeEncodingConfig(ModuleBasics), simapp.EmptyAppOptions{})
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
