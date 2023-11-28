package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"

	"github.com/Team-Kujira/core/x/cw-ica/types"
)

var (
	// TestOwnerAddress defines a reusable bech32 address for testing purposes
	TestOwnerAddress = "cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs"

	// TestAccountId defines a reusable interchainaccounts account id for testing purposes
	TestAccountId = "1"

	// TestMemo defines a reusable interchainaccounts memo for testing purposes
	TestMemo = "test memo"

	// TestVersion defines a reusable interchainaccounts version string for testing purposes
	TestVersion = string(icatypes.ModuleCdc.MustMarshalJSON(&icatypes.Metadata{
		Version:                icatypes.Version,
		ControllerConnectionId: ibctesting.FirstConnectionID,
		HostConnectionId:       ibctesting.FirstConnectionID,
		Encoding:               icatypes.EncodingProtobuf,
		TxType:                 icatypes.TxTypeSDKMultiMsg,
	}))

	TestMessage = &banktypes.MsgSend{
		FromAddress: "cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs",
		ToAddress:   "cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs",
		Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
	}
)

// TestMsgRegisterAccountValidateBasic tests ValidateBasic for MsgRegisterAccount
func TestMsgRegisterAccountValidateBasic(t *testing.T) {
	testCases := []struct {
		name    string
		msg     *types.MsgRegisterAccount
		expPass bool
	}{
		{"success", types.NewMsgRegisterAccount(TestOwnerAddress, ibctesting.FirstConnectionID, TestAccountId, TestVersion), true},
		{"account id is empty", types.NewMsgRegisterAccount(TestOwnerAddress, ibctesting.FirstConnectionID, "", TestVersion), false},
		{"owner address is empty", types.NewMsgRegisterAccount("", ibctesting.FirstConnectionID, TestAccountId, TestVersion), false},
		{"owner address is invalid", types.NewMsgRegisterAccount("invalid_address", ibctesting.FirstConnectionID, TestAccountId, TestVersion), false},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "valid test case %d failed: %s", i, tc.name)
		} else {
			require.Error(t, err, "invalid test case %d passed: %s", i, tc.name)
		}
	}
}

// TestMsgRegisterAccountGetSigners tests GetSigners for MsgRegisterAccount
func TestMsgRegisterAccountGetSigners(t *testing.T) {
	expSigner, err := sdk.AccAddressFromBech32(TestOwnerAddress)
	require.NoError(t, err)

	msg := types.NewMsgRegisterAccount(TestOwnerAddress, ibctesting.FirstConnectionID, TestAccountId, TestVersion)

	require.Equal(t, []sdk.AccAddress{expSigner}, msg.GetSigners())
}

// TestMsgSubmitTxValidateBasic tests ValidateBasic for MsgSubmitTx
func TestMsgSubmitTxValidateBasic(t *testing.T) {
	var msg *types.MsgSubmitTx

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"success",
			func() {},
			true,
		},
		{
			"owner address is invalid",
			func() {
				msg.Sender = "invalid_address"
			},
			false,
		},
	}

	for i, tc := range testCases {
		msg, _ = types.NewMsgSubmitTx(TestMessage, ibctesting.FirstConnectionID, TestAccountId, TestOwnerAddress, TestMemo, 1)

		tc.malleate()

		err := msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "valid test case %d failed: %s", i, tc.name)
		} else {
			require.Error(t, err, "invalid test case %d passed: %s", i, tc.name)
		}
	}
}

// TestMsgSubmitTxGetSigners tests GetSigners for MsgSubmitTx
func TestMsgSubmitTxGetSigners(t *testing.T) {
	expSigner, err := sdk.AccAddressFromBech32(TestOwnerAddress)
	require.NoError(t, err)

	msg, err := types.NewMsgSubmitTx(TestMessage, ibctesting.FirstConnectionID, TestAccountId, TestOwnerAddress, TestMemo, 1)
	require.NoError(t, err)

	require.Equal(t, []sdk.AccAddress{expSigner}, msg.GetSigners())
}
