package keeper_test

import (
	"github.com/Team-Kujira/core/x/onion/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (s *KeeperTestSuite) TestExecuteAnte() {
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

	specs := map[string]struct {
		msgs   []sdk.Msg
		accNum uint64
		nonce  uint64
		expErr bool
	}{
		"empty message tx ante": {
			msgs:   []sdk.Msg{},
			accNum: types.AccountNumber,
			nonce:  0,
			expErr: true,
		},
		"single message tx ante": {
			msgs:   []sdk.Msg{msgSend1},
			accNum: types.AccountNumber,
			nonce:  0,
			expErr: false,
		},
		"multiple messages tx ante": {
			msgs:   []sdk.Msg{msgSend1, msgSend2},
			accNum: types.AccountNumber,
			nonce:  0,
			expErr: false,
		},
		"invalid nonce check": {
			msgs:   []sdk.Msg{msgSend1},
			accNum: types.AccountNumber,
			nonce:  1,
			expErr: true,
		},
		"invalid account number check": {
			msgs:   []sdk.Msg{msgSend1},
			accNum: 1,
			nonce:  0,
			expErr: true,
		},
	}
	for msg, spec := range specs {
		spec := spec
		s.Run(msg, func() {
			s.SetupTest()
			s.Ctx = s.Ctx.WithChainID("test")
			tx := newTx(s.T(), s.App.TxConfig(), addr1, s.Ctx.ChainID(), spec.accNum, spec.msgs, spec.nonce, privKey1)
			err := s.App.OnionKeeper.ExecuteAnte(s.Ctx, tx)
			if spec.expErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				// check account set in account module
				acc1 := s.App.AccountKeeper.GetAccount(s.Ctx, addr1)
				s.Require().NotNil(acc1)
				// check sequence increase
				seq, err := s.App.OnionKeeper.GetSequence(s.Ctx, addr1.String())
				s.Require().NoError(err)
				s.Require().Equal(seq.Sequence, uint64(1))
			}
		})
	}
}
