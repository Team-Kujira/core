package wasmbinding_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/ed25519"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

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
