package keeper_test

import (
	"testing"
	"time"

	"github.com/Team-Kujira/core/app"
	tmtypes "github.com/cometbft/cometbft/proto/tendermint/types"
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
	suite.ctx = suite.app.BaseApp.NewContext(false, tmtypes.Header{Height: 1, ChainID: "kujira-1", Time: time.Now().UTC()})
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
