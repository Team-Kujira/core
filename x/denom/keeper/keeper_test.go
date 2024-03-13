package keeper_test

import (
	"testing"

	app "github.com/Team-Kujira/core/app"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	suite.Suite

	App *app.App
	Ctx sdk.Context
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	app := app.Setup(suite.T(), false)

	suite.Ctx = app.BaseApp.NewContext(false, tmproto.Header{})
	suite.App = app
}
