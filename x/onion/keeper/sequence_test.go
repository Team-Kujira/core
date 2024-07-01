package keeper_test

import (
	"github.com/Team-Kujira/core/x/onion/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestSequence() {
	suite.SetupTest()

	// Set accounts
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	addr3 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	err := suite.App.OnionKeeper.SetSequence(suite.Ctx, types.OnionSequence{
		Address:  addr1.String(),
		Sequence: 1,
	})
	suite.Require().NoError(err)
	err = suite.App.OnionKeeper.SetSequence(suite.Ctx, types.OnionSequence{
		Address:  addr2.String(),
		Sequence: 2,
	})
	suite.Require().NoError(err)

	// Check queries
	sequence, err := suite.App.OnionKeeper.GetSequence(suite.Ctx, addr1.String())
	suite.Require().NoError(err)
	suite.Require().Equal(sequence.Address, addr1.String())
	suite.Require().Equal(sequence.Sequence, uint64(1))
	sequence, err = suite.App.OnionKeeper.GetSequence(suite.Ctx, addr2.String())
	suite.Require().NoError(err)
	suite.Require().Equal(sequence.Address, addr2.String())
	suite.Require().Equal(sequence.Sequence, uint64(2))
	sequence, err = suite.App.OnionKeeper.GetSequence(suite.Ctx, addr3.String())
	suite.Require().NoError(err)
	suite.Require().Equal(sequence.Address, addr3.String())
	suite.Require().Equal(sequence.Sequence, uint64(0))
}
