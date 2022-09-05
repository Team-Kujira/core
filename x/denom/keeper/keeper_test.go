package keeper_test

// import (
// 	"testing"

// 	"github.com/Team-Kujira/core/x/denom/types"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/stretchr/testify/suite"
// )

// type KeeperTestSuite struct {
// 	apptesting.KeeperTestHelper

// 	queryClient types.QueryClient
// }

// func TestKeeperTestSuite(t *testing.T) {
// 	suite.Run(t, new(KeeperTestSuite))
// }

// func (suite *KeeperTestSuite) SetupTest() {
// 	suite.Setup()

// 	// Fund every TestAcc with 100 denom creation fees.
// 	fundAccsAmount := sdk.NewCoins(sdk.NewCoin(types.DefaultParams().CreationFee[0].Denom, types.DefaultParams().CreationFee[0].Amount.MulRaw(100)))
// 	for _, acc := range suite.TestAccs {
// 		suite.FundAcc(acc, fundAccsAmount)
// 	}

// 	suite.Setupdenom()

// 	suite.queryClient = types.NewQueryClient(suite.QueryHelper)
// }
