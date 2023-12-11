package app

import (
	"encoding/json"
	"os"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	abci "github.com/cometbft/cometbft/abci/types"
	cmtjson "github.com/cometbft/cometbft/libs/json"
	cmttypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
)

// Setup initializes a new KujiraApp.
func Setup(t *testing.T, isCheckTx bool) *App {
	db := dbm.NewMemDB()
	var wasmOpts []wasmkeeper.Option
	appOptions := make(simtestutil.AppOptionsMap, 0)

	app := New(
		log.NewNopLogger(),
		db,
		nil,
		true,
		appOptions,
		wasmOpts,
	)

	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err)
	// create validator set with single validator
	validator := cmttypes.NewValidator(pubKey, 1)
	valSet := cmttypes.NewValidatorSet([]*cmttypes.Validator{validator})

	// generate genesis account
	senderPrivKey := secp256k1.GenPrivKey()
	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100000000000000))),
	}

	if !isCheckTx {
		genesisState := app.DefaultGenesis()
		genesisState, err = simtestutil.GenesisStateWithValSet(
			app.AppCodec(),
			genesisState,
			valSet,
			[]authtypes.GenesisAccount{acc},
			balance,
		)
		if err != nil {
			panic(err)
		}

		stateBytes, err := cmtjson.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		_, err = app.InitChain(
			&abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simtestutil.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
		require.NoError(t, err)
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
	var wasmOpts []wasmkeeper.Option
	appOptions := make(simtestutil.AppOptionsMap, 0)

	app = New(
		log.NewNopLogger(),
		db,
		nil,
		true,
		appOptions,
		wasmOpts,
	)
	if !isCheckTx {
		genesisState := app.DefaultGenesis()
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		_, err = app.InitChain(
			&abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simtestutil.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
		if err != nil {
			panic(err)
		}
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
