package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/require"
)

func (s *KeeperTestSuite) TestExecuteTxMsgs() {
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
		"empty messages execution": {
			msgs:               []sdk.Msg{},
			expErr:             false,
			expSenderBalance:   sdk.Coins{},
			expReceiverBalance: sdk.Coins{},
		},
		"successful execution of a single message": {
			msgs:               []sdk.Msg{msgSend1},
			expErr:             false,
			expSenderBalance:   sdk.Coins{},
			expReceiverBalance: sdk.Coins{},
		},
		"successful execution of multiple messages": {
			msgs:               []sdk.Msg{msgSend1, msgSend2},
			expErr:             false,
			expSenderBalance:   sdk.Coins{},
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

			tx := newTx(s.T(), s.App.TxConfig(), spec.msgs)
			results, err := s.App.OnionKeeper.ExecuteTxMsgs(s.Ctx, tx)
			if spec.expErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Len(results, len(spec.msgs))
			}
		})
	}
}

func newTx(t *testing.T, cfg client.TxConfig, msgs []sdk.Msg) signing.Tx {
	builder := cfg.NewTxBuilder()
	builder.SetMsgs(msgs...)
	nonce := uint64(1)
	setTxSignature(t, builder, nonce)

	return builder.GetTx()
}

func setTxSignature(t *testing.T, builder client.TxBuilder, nonce uint64) {
	privKey := secp256k1.GenPrivKeyFromSecret([]byte("test"))
	pubKey := privKey.PubKey()
	err := builder.SetSignatures(
		signingtypes.SignatureV2{
			PubKey:   pubKey,
			Sequence: nonce,
			Data:     &signingtypes.SingleSignatureData{},
		},
	)
	require.NoError(t, err)
}
