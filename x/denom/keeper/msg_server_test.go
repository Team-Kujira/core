package keeper_test

import (
	"github.com/Team-Kujira/core/x/denom/keeper"
	"github.com/Team-Kujira/core/x/denom/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

func (suite *KeeperTestSuite) TestMsgServerCreateDenom() {
	for _, tc := range []struct {
		desc         string
		balance      sdk.Coins
		whitelisted  bool
		balanceAfter sdk.Coins
		expPass      bool
	}{
		{
			desc:         "successful denom creation: positive fee payment",
			balance:      sdk.Coins{sdk.NewInt64Coin("ukuji", 10_000_000)},
			whitelisted:  false,
			balanceAfter: sdk.Coins{},
			expPass:      true,
		},
		{
			desc:         "insufficient balance for fee payment",
			balance:      sdk.Coins{},
			whitelisted:  false,
			balanceAfter: sdk.Coins{},
			expPass:      false,
		},
		{
			desc:         "successful denom creation: whitelisted address",
			balance:      sdk.Coins{},
			whitelisted:  true,
			balanceAfter: sdk.Coins{},
			expPass:      true,
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()

			// bootstrap accounts
			sender := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

			if tc.whitelisted {
				suite.App.DenomKeeper.SetNoFeeAccount(suite.Ctx, sender.String())
			}

			if tc.balance.IsAllPositive() {
				err := suite.App.BankKeeper.MintCoins(suite.Ctx, minttypes.ModuleName, tc.balance)
				suite.Require().NoError(err)
				err = suite.App.BankKeeper.SendCoinsFromModuleToAccount(suite.Ctx, minttypes.ModuleName, sender, tc.balance)
				suite.Require().NoError(err)
			}

			msgServer := keeper.NewMsgServerImpl(*suite.App.DenomKeeper)
			resp, err := msgServer.CreateDenom(
				sdk.WrapSDKContext(suite.Ctx),
				&types.MsgCreateDenom{
					Sender: sender.String(),
					Nonce:  "1",
				},
			)
			if !tc.expPass {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				denom, err := types.GetTokenDenom(sender.String(), "1")
				suite.Require().NoError(err)
				suite.Require().Equal(resp.NewTokenDenom, denom)

				// check balance after
				balances := suite.App.BankKeeper.GetAllBalances(suite.Ctx, sender)
				suite.Require().Equal(balances.String(), tc.balanceAfter.String())

				metadata, err := suite.App.DenomKeeper.GetAuthorityMetadata(suite.Ctx, denom)
				suite.Require().NoError(err)
				suite.Require().Equal(metadata.Admin, sender.String())
			}
		})
	}
}
