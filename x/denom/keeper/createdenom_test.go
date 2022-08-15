package keeper_test

// import (
// 	sdk "github.com/cosmos/cosmos-sdk/types"

// 	"github.com/Team-Kujira/core/x/denom/keeper"
// 	"github.com/Team-Kujira/core/x/denom/types"
// )

// func (suite *KeeperTestSuite) TestMsgCreateDenom() {
// 	suite.SetupTest()

// 	msgServer := keeper.NewMsgServerImpl(*suite.App.denomKeeper)

// 	creationFee := suite.App.denomKeeper.GetParams(suite.Ctx).CreationFee

// 	// Get balance of acc 0 before creating a denom
// 	preCreateBalance := suite.App.BankKeeper.GetBalance(suite.Ctx, suite.TestAccs[0], creationFee[0].Denom)

// 	// Creating a denom should work
// 	res, err := msgServer.CreateDenom(sdk.WrapSDKContext(suite.Ctx), types.NewMsgCreateDenom(suite.TestAccs[0].String(), "bitcoin"))
// 	suite.Require().NoError(err)
// 	suite.Require().NotEmpty(res.GetNewTokenDenom())

// 	// Make sure that the admin is set correctly
// 	queryRes, err := suite.queryClient.DenomAuthorityMetadata(suite.Ctx.Context(), &types.QueryDenomAuthorityMetadataRequest{
// 		Denom: res.GetNewTokenDenom(),
// 	})
// 	suite.Require().NoError(err)
// 	suite.Require().Equal(suite.TestAccs[0].String(), queryRes.AuthorityMetadata.Admin)

// 	// Make sure that creation fee was deducted
// 	postCreateBalance := suite.App.BankKeeper.GetBalance(suite.Ctx, suite.TestAccs[0], suite.App.denomKeeper.GetParams(suite.Ctx).CreationFee[0].Denom)
// 	suite.Require().True(preCreateBalance.Sub(postCreateBalance).IsEqual(creationFee[0]))

// 	// Make sure that a second version of the same denom can't be recreated
// 	res, err = msgServer.CreateDenom(sdk.WrapSDKContext(suite.Ctx), types.NewMsgCreateDenom(suite.TestAccs[0].String(), "bitcoin"))
// 	suite.Require().Error(err)

// 	// Creating a second denom should work
// 	res, err = msgServer.CreateDenom(sdk.WrapSDKContext(suite.Ctx), types.NewMsgCreateDenom(suite.TestAccs[0].String(), "litecoin"))
// 	suite.Require().NoError(err)
// 	suite.Require().NotEmpty(res.GetNewTokenDenom())

// 	// Try querying all the denoms created by suite.TestAccs[0]
// 	queryRes2, err := suite.queryClient.DenomsFromCreator(suite.Ctx.Context(), &types.QueryDenomsFromCreatorRequest{
// 		Creator: suite.TestAccs[0].String(),
// 	})
// 	suite.Require().NoError(err)
// 	suite.Require().Len(queryRes2.Denoms, 2)

// 	// Make sure that a second account can create a denom with the same nonce
// 	res, err = msgServer.CreateDenom(sdk.WrapSDKContext(suite.Ctx), types.NewMsgCreateDenom(suite.TestAccs[1].String(), "bitcoin"))
// 	suite.Require().NoError(err)
// 	suite.Require().NotEmpty(res.GetNewTokenDenom())

// 	// Make sure that an address with a "/" in it can't create denoms
// 	res, err = msgServer.CreateDenom(sdk.WrapSDKContext(suite.Ctx), types.NewMsgCreateDenom("addr.eth/creator", "bitcoin"))
// 	suite.Require().Error(err)
// }
