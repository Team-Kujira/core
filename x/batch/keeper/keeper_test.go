package keeper_test

import (
	"testing"

	"github.com/Team-Kujira/core/app"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context

	app *app.App
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = app.Setup(suite.T(), false)
	suite.ctx = suite.app.BaseApp.NewContext(false)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
