package wasmbinding_test

import (
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Team-Kujira/core/x/denom/wasm"

	"github.com/Team-Kujira/core/x/denom/types"

	"github.com/stretchr/testify/require"

	"github.com/Team-Kujira/core/app"

	"github.com/cosmos/cosmos-sdk/x/bank/testutil"

	tmtypes "github.com/cometbft/cometbft/proto/tendermint/types"
)

func fundAccount(t *testing.T, ctx sdk.Context, app *app.App, addr sdk.AccAddress, coins sdk.Coins) {
	err := testutil.FundAccount(
		ctx,
		app.BankKeeper,
		addr,
		coins,
	)
	require.NoError(t, err)
}

func TestCreateDenom(t *testing.T) {
	actor := RandomAccountAddress()
	app := app.Setup(t, false)
	ctx := app.BaseApp.NewContextLegacy(false, tmtypes.Header{Height: 1, ChainID: "kujira-1", Time: time.Now().UTC()})

	// Fund actor with 100 base denom creation feesme
	actorAmount := sdk.NewCoins(sdk.NewCoin(types.DefaultParams().CreationFee[0].Denom, types.DefaultParams().CreationFee[0].Amount.MulRaw(100)))

	fundAccount(t, ctx, app, actor, actorAmount)

	specs := map[string]struct {
		createDenom *wasm.Create
		expErr      bool
	}{
		"valid sub-denom": {
			createDenom: &wasm.Create{
				Subdenom: "MOON",
			},
		},
		"empty sub-denom": {
			createDenom: &wasm.Create{
				Subdenom: "",
			},
			expErr: false,
		},
		"invalid sub-denom": {
			createDenom: &wasm.Create{
				Subdenom: "adsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsf",
			},
			expErr: true,
		},
		"null create denom": {
			createDenom: nil,
			expErr:      true,
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			// when
			gotErr := wasm.PerformCreate(*app.DenomKeeper, app.BankKeeper, ctx, actor, spec.createDenom)
			// then
			if spec.expErr {
				require.Error(t, gotErr)
				return
			}
			require.NoError(t, gotErr)
		})
	}
}

func TestChangeAdmin(t *testing.T) {
	const validDenom = "validdenom"

	tokenCreator := RandomAccountAddress()

	specs := map[string]struct {
		actor       sdk.AccAddress
		changeAdmin *wasm.ChangeAdmin

		expErrMsg string
	}{
		"valid": {
			changeAdmin: &wasm.ChangeAdmin{
				Denom:   fmt.Sprintf("factory/%s/%s", tokenCreator.String(), validDenom),
				Address: RandomBech32AccountAddress(),
			},
			actor: tokenCreator,
		},
		"typo in factory in denom name": {
			changeAdmin: &wasm.ChangeAdmin{
				Denom:   fmt.Sprintf("facory/%s/%s", tokenCreator.String(), validDenom),
				Address: RandomBech32AccountAddress(),
			},
			actor:     tokenCreator,
			expErrMsg: "denom prefix is incorrect. Is: facory.  Should be: factory: invalid denom",
		},
		"invalid address in denom": {
			changeAdmin: &wasm.ChangeAdmin{
				Denom:   fmt.Sprintf("factory/%s/%s", RandomBech32AccountAddress(), validDenom),
				Address: RandomBech32AccountAddress(),
			},
			actor:     tokenCreator,
			expErrMsg: "failed changing admin from message: unauthorized account",
		},
		"other denom name in 3 part name": {
			changeAdmin: &wasm.ChangeAdmin{
				Denom:   fmt.Sprintf("factory/%s/%s", tokenCreator.String(), "invalid denom"),
				Address: RandomBech32AccountAddress(),
			},
			actor:     tokenCreator,
			expErrMsg: fmt.Sprintf("invalid denom: factory/%s/invalid denom", tokenCreator.String()),
		},
		"empty denom": {
			changeAdmin: &wasm.ChangeAdmin{
				Denom:   "",
				Address: RandomBech32AccountAddress(),
			},
			actor:     tokenCreator,
			expErrMsg: "invalid denom: ",
		},
		"empty address": {
			changeAdmin: &wasm.ChangeAdmin{
				Denom:   fmt.Sprintf("factory/%s/%s", tokenCreator.String(), validDenom),
				Address: "",
			},
			actor:     tokenCreator,
			expErrMsg: "address from bech32: empty address string is not allowed",
		},
		"creator is a different address": {
			changeAdmin: &wasm.ChangeAdmin{
				Denom:   fmt.Sprintf("factory/%s/%s", tokenCreator.String(), validDenom),
				Address: RandomBech32AccountAddress(),
			},
			actor:     RandomAccountAddress(),
			expErrMsg: "failed changing admin from message: unauthorized account",
		},
		"change to the same address": {
			changeAdmin: &wasm.ChangeAdmin{
				Denom:   fmt.Sprintf("factory/%s/%s", tokenCreator.String(), validDenom),
				Address: tokenCreator.String(),
			},
			actor: tokenCreator,
		},
		"nil binding": {
			actor:     tokenCreator,
			expErrMsg: "invalid request: changeAdmin is nil - original request: ",
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			// Setup
			app := app.Setup(t, false)
			ctx := app.BaseApp.NewContextLegacy(false, tmtypes.Header{Height: 1, ChainID: "kujira-1", Time: time.Now().UTC()})

			// Fund actor with 100 base denom creation fees
			actorAmount := sdk.NewCoins(sdk.NewCoin(types.DefaultParams().CreationFee[0].Denom, types.DefaultParams().CreationFee[0].Amount.MulRaw(100)))
			fundAccount(t, ctx, app, tokenCreator, actorAmount)

			err := wasm.PerformCreate(*app.DenomKeeper, app.BankKeeper, ctx, tokenCreator, &wasm.Create{
				Subdenom: validDenom,
			})
			require.NoError(t, err)

			err = wasm.PerformChangeAdmin(*app.DenomKeeper, ctx, spec.actor, spec.changeAdmin)
			if len(spec.expErrMsg) > 0 {
				require.Error(t, err)
				actualErrMsg := err.Error()
				require.Equal(t, spec.expErrMsg, actualErrMsg)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMint(t *testing.T) {
	creator := RandomAccountAddress()
	app := app.Setup(t, false)
	ctx := app.BaseApp.NewContextLegacy(false, tmtypes.Header{Height: 1, ChainID: "kujira-1", Time: time.Now().UTC()})

	// Fund actor with 100 base denom creation fees
	tokenCreationFeeAmt := sdk.NewCoins(sdk.NewCoin(types.DefaultParams().CreationFee[0].Denom, types.DefaultParams().CreationFee[0].Amount.MulRaw(100)))
	fundAccount(t, ctx, app, creator, tokenCreationFeeAmt)

	// Create denoms for valid mint tests
	validDenom := wasm.Create{
		Subdenom: "MOON",
	}
	err := wasm.PerformCreate(*app.DenomKeeper, app.BankKeeper, ctx, creator, &validDenom)
	require.NoError(t, err)

	emptyDenom := wasm.Create{
		Subdenom: "",
	}
	err = wasm.PerformCreate(*app.DenomKeeper, app.BankKeeper, ctx, creator, &emptyDenom)
	require.NoError(t, err)

	validDenomStr := fmt.Sprintf("factory/%s/%s", creator.String(), validDenom.Subdenom)
	emptyDenomStr := fmt.Sprintf("factory/%s/%s", creator.String(), emptyDenom.Subdenom)

	lucky := RandomAccountAddress()

	// lucky was broke
	balances := app.BankKeeper.GetAllBalances(ctx, lucky)
	require.Empty(t, balances)

	amount, ok := math.NewIntFromString("8080")
	require.True(t, ok)

	specs := map[string]struct {
		mint   *wasm.Mint
		expErr bool
	}{
		"valid mint": {
			mint: &wasm.Mint{
				Denom:     validDenomStr,
				Amount:    amount,
				Recipient: lucky.String(),
			},
		},
		"empty sub-denom": {
			mint: &wasm.Mint{
				Denom:     emptyDenomStr,
				Amount:    amount,
				Recipient: lucky.String(),
			},
			expErr: false,
		},
		"nonexistent sub-denom": {
			mint: &wasm.Mint{
				Denom:     fmt.Sprintf("factory/%s/%s", creator.String(), "SUN"),
				Amount:    amount,
				Recipient: lucky.String(),
			},
			expErr: true,
		},
		"invalid sub-denom": {
			mint: &wasm.Mint{
				Denom:     "sub-denom_2",
				Amount:    amount,
				Recipient: lucky.String(),
			},
			expErr: true,
		},
		"zero amount": {
			mint: &wasm.Mint{
				Denom:     validDenomStr,
				Amount:    math.ZeroInt(),
				Recipient: lucky.String(),
			},
			expErr: false,
		},
		"negative amount": {
			mint: &wasm.Mint{
				Denom:     validDenomStr,
				Amount:    amount.Neg(),
				Recipient: lucky.String(),
			},
			expErr: true,
		},
		"empty recipient": {
			mint: &wasm.Mint{
				Denom:     validDenomStr,
				Amount:    amount,
				Recipient: "",
			},
			expErr: true,
		},
		"invalid recipient": {
			mint: &wasm.Mint{
				Denom:     validDenomStr,
				Amount:    amount,
				Recipient: "invalid",
			},
			expErr: true,
		},
		"null mint": {
			mint:   nil,
			expErr: true,
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			// when
			gotErr := wasm.PerformMint(*app.DenomKeeper, app.BankKeeper, ctx, creator, spec.mint)
			// then
			if spec.expErr {
				require.Error(t, gotErr)
				return
			}
			require.NoError(t, gotErr)
		})
	}
}

func TestBurn(t *testing.T) {
	creator := RandomAccountAddress()
	app := app.Setup(t, false)
	ctx := app.BaseApp.NewContextLegacy(false, tmtypes.Header{Height: 1, ChainID: "kujira-1", Time: time.Now().UTC()})

	// Fund actor with 100 base denom creation fees
	tokenCreationFeeAmt := sdk.NewCoins(sdk.NewCoin(types.DefaultParams().CreationFee[0].Denom, types.DefaultParams().CreationFee[0].Amount.MulRaw(100)))
	fundAccount(t, ctx, app, creator, tokenCreationFeeAmt)

	// Create denoms for valid burn tests
	validDenom := wasm.Create{
		Subdenom: "MOON",
	}
	err := wasm.PerformCreate(*app.DenomKeeper, app.BankKeeper, ctx, creator, &validDenom)
	require.NoError(t, err)

	emptyDenom := wasm.Create{
		Subdenom: "",
	}
	err = wasm.PerformCreate(*app.DenomKeeper, app.BankKeeper, ctx, creator, &emptyDenom)
	require.NoError(t, err)

	lucky := RandomAccountAddress()

	// lucky was broke
	balances := app.BankKeeper.GetAllBalances(ctx, lucky)
	require.Empty(t, balances)

	validDenomStr := fmt.Sprintf("factory/%s/%s", creator.String(), validDenom.Subdenom)
	emptyDenomStr := fmt.Sprintf("factory/%s/%s", creator.String(), emptyDenom.Subdenom)
	mintAmount, ok := math.NewIntFromString("8080")
	require.True(t, ok)

	specs := map[string]struct {
		sender sdk.AccAddress
		burn   *wasm.Burn
		expErr bool
	}{
		"valid burn": {
			sender: creator,
			burn: &wasm.Burn{
				Denom:  validDenomStr,
				Amount: mintAmount,
			},
			expErr: false,
		},
		"non admin address": {
			sender: lucky,
			burn: &wasm.Burn{
				Denom:  validDenomStr,
				Amount: mintAmount,
			},
			expErr: true,
		},
		"empty sub-denom": {
			sender: creator,
			burn: &wasm.Burn{
				Denom:  emptyDenomStr,
				Amount: mintAmount,
			},
			expErr: false,
		},
		"invalid sub-denom": {
			sender: creator,
			burn: &wasm.Burn{
				Denom:  "sub-denom_2",
				Amount: mintAmount,
			},
			expErr: true,
		},
		"non-minted denom": {
			sender: creator,
			burn: &wasm.Burn{
				Denom:  fmt.Sprintf("factory/%s/%s", creator.String(), "SUN"),
				Amount: mintAmount,
			},
			expErr: true,
		},
		"zero amount": {
			sender: creator,

			burn: &wasm.Burn{
				Denom:  validDenomStr,
				Amount: math.ZeroInt(),
			},
			expErr: false,
		},
		"negative amount": {
			sender: creator,

			burn:   nil,
			expErr: true,
		},
		"null burn": {
			sender: creator,

			burn: &wasm.Burn{
				Denom:  validDenomStr,
				Amount: mintAmount.Neg(),
			},
			expErr: true,
		},
	}

	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			// Mint valid denom str and empty denom string for burn test
			mintBinding := &wasm.Mint{
				Denom:     validDenomStr,
				Amount:    mintAmount,
				Recipient: creator.String(),
			}
			err := wasm.PerformMint(*app.DenomKeeper, app.BankKeeper, ctx, creator, mintBinding)
			require.NoError(t, err)

			emptyDenomMintBinding := &wasm.Mint{
				Denom:     emptyDenomStr,
				Amount:    mintAmount,
				Recipient: creator.String(),
			}
			err = wasm.PerformMint(*app.DenomKeeper, app.BankKeeper, ctx, creator, emptyDenomMintBinding)
			require.NoError(t, err)

			// when
			gotErr := wasm.PerformBurn(*app.DenomKeeper, ctx, spec.sender, spec.burn)
			// then
			if spec.expErr {
				require.Error(t, gotErr)
				return
			}
			require.NoError(t, gotErr)
		})
	}
}
