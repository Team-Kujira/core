package keeper_test

import (
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestNoFeeAccounts() {
	suite.SetupTest()

	// Set accounts
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	addr3 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	suite.App.DenomKeeper.SetNoFeeAccount(suite.Ctx, addr1.String())
	suite.App.DenomKeeper.SetNoFeeAccount(suite.Ctx, addr2.String())

	// Check queries
	suite.Require().True(suite.App.DenomKeeper.IsNoFeeAccount(suite.Ctx, addr1.String()))
	suite.Require().True(suite.App.DenomKeeper.IsNoFeeAccount(suite.Ctx, addr2.String()))
	suite.Require().False(suite.App.DenomKeeper.IsNoFeeAccount(suite.Ctx, addr3.String()))
	suite.Require().Len(suite.App.DenomKeeper.GetNoFeeAccounts(suite.Ctx), 2)

	// Remove accounts
	suite.App.DenomKeeper.RemoveNoFeeAccount(suite.Ctx, addr2.String())
	suite.App.DenomKeeper.RemoveNoFeeAccount(suite.Ctx, addr3.String())

	// Check queries
	suite.Require().True(suite.App.DenomKeeper.IsNoFeeAccount(suite.Ctx, addr1.String()))
	suite.Require().False(suite.App.DenomKeeper.IsNoFeeAccount(suite.Ctx, addr2.String()))
	suite.Require().False(suite.App.DenomKeeper.IsNoFeeAccount(suite.Ctx, addr3.String()))
	suite.Require().Len(suite.App.DenomKeeper.GetNoFeeAccounts(suite.Ctx), 1)
}
