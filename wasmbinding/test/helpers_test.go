package wasmbinding_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/ed25519"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"

	"github.com/Team-Kujira/core/app"
)

func CreateTestInput(t *testing.T) (*app.App, sdk.Context) {
	app := app.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "kujira-1", Time: time.Now().UTC()})
	return app, ctx
}

func FundAccount(t *testing.T, ctx sdk.Context, app *app.App, acct sdk.AccAddress) {
	err := testutil.FundAccount(app.BankKeeper, ctx, acct, sdk.NewCoins(
		sdk.NewCoin("uosmo", sdk.NewInt(10000000000)),
	))
	require.NoError(t, err)
}

// SetupContract uploads the wasm code to the wasm storage and creates a contract instance
func SetupContract(t *testing.T, ctx sdk.Context, app *app.App) (sdk.AccAddress, sdk.AccAddress) {
	// create an addr that perform contract transactions, upload the wasm code
	actor := RandomAccountAddress()
	storeWasmCode(t, ctx, app, actor)

	cInfo := app.WasmKeeper.GetCodeInfo(ctx, 1)
	require.NotNil(t, cInfo)

	// create a contract instance
	contractAddr := instantiateContract(t, ctx, app, actor)
	require.NotEmpty(t, contractAddr)

	return contractAddr, actor
}

func storeWasmCode(t *testing.T, ctx sdk.Context, app *app.App, addr sdk.AccAddress) {
	wasmCode, err := os.ReadFile("../testdata/kujira_ibc.wasm")
	require.NoError(t, err)

	contractKeeper := keeper.NewDefaultPermissionKeeper(app.WasmKeeper)
	_, _, err = contractKeeper.Create(ctx, addr, wasmCode, &wasmtypes.AccessConfig{
		Permission: wasmtypes.AccessTypeEverybody,
	})
	require.NoError(t, err)
}

func instantiateContract(t *testing.T, ctx sdk.Context, app *app.App, funder sdk.AccAddress) sdk.AccAddress {
	initMsgBz := []byte("{}")
	contractKeeper := keeper.NewDefaultPermissionKeeper(app.WasmKeeper)
	codeID := uint64(1)
	addr, _, err := contractKeeper.Instantiate(ctx, codeID, funder, funder, initMsgBz, "demo contract", nil)
	require.NoError(t, err)

	return addr
}

// we need to make this deterministic (same every test run), as content might affect gas costs
func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

func RandomAccountAddress() sdk.AccAddress {
	_, _, addr := keyPubAddr()
	return addr
}

func RandomBech32AccountAddress() string {
	return RandomAccountAddress().String()
}
