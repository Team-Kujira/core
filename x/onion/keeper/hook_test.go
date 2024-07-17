package keeper_test

import (
	"encoding/base64"

	"github.com/Team-Kujira/core/x/onion/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

func (s *KeeperTestSuite) TestOnReceivePacketHook() {
	privKey1 := secp256k1.GenPrivKeyFromSecret([]byte("test1"))
	pubKey1 := privKey1.PubKey()
	privKey2 := secp256k1.GenPrivKeyFromSecret([]byte("test2"))
	pubKey2 := privKey2.PubKey()

	addr1 := sdk.AccAddress(pubKey1.Address())
	addr2 := sdk.AccAddress(pubKey2.Address())

	msgSend1 := &banktypes.MsgSend{
		FromAddress: addr1.String(),
		ToAddress:   addr2.String(),
		Amount:      sdk.Coins{sdk.NewInt64Coin("test", 100)},
	}
	msgSend2 := &banktypes.MsgSend{
		FromAddress: addr1.String(),
		ToAddress:   addr2.String(),
		Amount:      sdk.Coins{sdk.NewInt64Coin("test", 200)},
	}
	msgSend3 := &banktypes.MsgSend{
		FromAddress: addr1.String(),
		ToAddress:   addr2.String(),
		Amount:      sdk.Coins{sdk.NewInt64Coin("test", 1000)},
	}

	specs := map[string]struct {
		msgs               []sdk.Msg
		expErr             bool
		expSenderBalance   sdk.Coins
		expReceiverBalance sdk.Coins
	}{
		"successful execution of a single message": {
			msgs:               []sdk.Msg{msgSend1},
			expErr:             false,
			expSenderBalance:   sdk.Coins{sdk.NewInt64Coin("test", 400)},
			expReceiverBalance: sdk.Coins{sdk.NewInt64Coin("test", 100)},
		},
		"successful execution of multiple messages": {
			msgs:               []sdk.Msg{msgSend1, msgSend2},
			expErr:             false,
			expSenderBalance:   sdk.Coins{sdk.NewInt64Coin("test", 200)},
			expReceiverBalance: sdk.Coins{sdk.NewInt64Coin("test", 300)},
		},
		"empty messages execution": {
			msgs:               []sdk.Msg{},
			expErr:             true,
			expSenderBalance:   sdk.Coins{sdk.NewInt64Coin("test", 500)},
			expReceiverBalance: sdk.Coins{},
		},
		"one execution failure in multiple messages": {
			msgs:               []sdk.Msg{msgSend1, msgSend3},
			expErr:             true,
			expSenderBalance:   sdk.Coins{},
			expReceiverBalance: sdk.Coins{},
		},
	}
	for msg, spec := range specs {
		spec := spec
		s.Run(msg, func() {
			s.SetupTest()
			coins := sdk.Coins{sdk.NewInt64Coin("test", 500)}
			err := s.App.BankKeeper.MintCoins(s.Ctx, minttypes.ModuleName, coins)
			s.Require().NoError(err)
			err = s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, minttypes.ModuleName, addr1, coins)
			s.Require().NoError(err)

			tx := newTx(s.T(), s.App.TxConfig(), addr1, s.Ctx.ChainID(), types.AccountNumber, spec.msgs, 0, privKey1)
			txBytes, err := s.App.TxConfig().TxEncoder()(tx)
			s.Require().NoError(err)

			memo := base64.StdEncoding.EncodeToString(txBytes)

			s.App.OnionKeeper.HandleTransferHook(s.Ctx, memo, s.App.TxConfig())
			if spec.expErr {
				// Check sequence change
				seq, err := s.App.OnionKeeper.GetSequence(s.Ctx, addr1.String())
				s.Require().NoError(err)
				s.Require().Equal(seq.Sequence, uint64(0))
			} else {
				// Check sequence change
				seq, err := s.App.OnionKeeper.GetSequence(s.Ctx, addr1.String())
				s.Require().NoError(err)
				s.Require().Equal(seq.Sequence, uint64(1))

				// Check balance
				senderBalance := s.App.BankKeeper.GetAllBalances(s.Ctx, addr1)
				s.Require().Equal(senderBalance.String(), spec.expSenderBalance.String())
				receiverBalance := s.App.BankKeeper.GetAllBalances(s.Ctx, addr2)
				s.Require().Equal(receiverBalance.String(), spec.expReceiverBalance.String())
			}
		})
	}
}
