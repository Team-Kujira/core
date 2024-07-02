package keeper

import (
	"bytes"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/Team-Kujira/core/x/onion/types"
	kmultisig "github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	txsigning "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func GetSignerAcc(ctx sdk.Context, ak types.AccountKeeper, addr sdk.AccAddress) (authtypes.AccountI, error) {
	if acc := ak.GetAccount(ctx, addr); acc != nil {
		return acc, nil
	}

	return nil, errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", addr)
}

// CountSubKeys counts the total number of keys for a multi-sig public key.
func CountSubKeys(pub cryptotypes.PubKey) int {
	v, ok := pub.(*kmultisig.LegacyAminoPubKey)
	if !ok {
		return 1
	}

	numKeys := 0
	for _, subkey := range v.GetPubKeys() {
		numKeys += CountSubKeys(subkey)
	}

	return numKeys
}

func OnlyLegacyAminoSigners(sigData txsigning.SignatureData) bool {
	switch v := sigData.(type) {
	case *txsigning.SingleSignatureData:
		return v.SignMode == txsigning.SignMode_SIGN_MODE_LEGACY_AMINO_JSON
	case *txsigning.MultiSignatureData:
		for _, s := range v.Signatures {
			if !OnlyLegacyAminoSigners(s) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func (k Keeper) ExecuteAnte(ctx sdk.Context, tx sdk.Tx) error {
	// ValidateBasicDecorator
	if err := tx.ValidateBasic(); err != nil {
		return err
	}

	// SetPubKeyDecorator
	sigTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return errorsmod.Wrap(sdkerrors.ErrTxDecode, "invalid tx type")
	}

	pubkeys, err := sigTx.GetPubKeys()
	if err != nil {
		return err
	}
	signers := sigTx.GetSigners()

	for i, pk := range pubkeys {
		if pk == nil {
			continue
		}
		if !bytes.Equal(pk.Address(), signers[i]) {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidPubKey,
				"pubKey does not match signer address %s with signer index: %d", signers[i], i)
		}

		acc, err := GetSignerAcc(ctx, k.accountKeeper, signers[i])
		if err != nil {
			return err
		}
		if acc.GetPubKey() != nil {
			continue
		}
		err = acc.SetPubKey(pk)
		if err != nil {
			return errorsmod.Wrap(sdkerrors.ErrInvalidPubKey, err.Error())
		}
		k.accountKeeper.SetAccount(ctx, acc)
	}

	// ValidateSigCountDecorator
	params := k.accountKeeper.GetParams(ctx)
	pubKeys, err := sigTx.GetPubKeys()
	if err != nil {
		return err
	}

	sigCount := 0
	for _, pk := range pubKeys {
		sigCount += CountSubKeys(pk)
		if uint64(sigCount) > params.TxSigLimit {
			return errorsmod.Wrapf(sdkerrors.ErrTooManySignatures,
				"signatures: %d, limit: %d", sigCount, params.TxSigLimit)
		}
	}

	// SigVerificationDecorator
	sigs, err := sigTx.GetSignaturesV2()
	if err != nil {
		return err
	}

	signerAddrs := sigTx.GetSigners()

	if len(sigs) != len(signerAddrs) {
		return errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "invalid number of signer;  expected: %d, got %d", len(signerAddrs), len(sigs))
	}

	for i, sig := range sigs {
		acc, err := GetSignerAcc(ctx, k.accountKeeper, signerAddrs[i])
		if err != nil {
			return err
		}

		pubKey := acc.GetPubKey()
		if pubKey == nil {
			return errorsmod.Wrap(sdkerrors.ErrInvalidPubKey, "pubkey on account is not set")
		}

		onionSeq := uint64(0)
		seq, err := k.GetSequence(ctx, acc.GetAddress().String())
		if err == nil {
			onionSeq = seq.Sequence
		}

		if sig.Sequence != onionSeq {
			return errorsmod.Wrapf(
				sdkerrors.ErrWrongSequence,
				"onion sequence mismatch, expected %d, got %d", onionSeq, sig.Sequence,
			)
		}

		chainID := ctx.ChainID()
		accNum := types.AccountNumber
		signerData := authsigning.SignerData{
			Address:       acc.GetAddress().String(),
			ChainID:       chainID,
			AccountNumber: accNum,
			Sequence:      onionSeq,
			PubKey:        pubKey,
		}

		err = authsigning.VerifySignature(pubKey, signerData, sig.Data, k.signModeHandler, tx)
		if err != nil {
			var errMsg string
			if OnlyLegacyAminoSigners(sig.Data) {
				errMsg = fmt.Sprintf("signature verification failed; please verify account number (%d), sequence (%d) and chain-id (%s)", accNum, acc.GetSequence(), chainID)
			} else {
				errMsg = fmt.Sprintf("signature verification failed; please verify account number (%d) and chain-id (%s)", accNum, chainID)
			}
			return errorsmod.Wrap(sdkerrors.ErrUnauthorized, errMsg)
		}
	}

	// IncrementSequenceDecorator
	for _, addr := range sigTx.GetSigners() {
		seq, err := k.GetSequence(ctx, addr.String())
		if err != nil {
			seq = types.OnionSequence{
				Address:  addr.String(),
				Sequence: 0,
			}
		}

		seq.Sequence++
		err = k.SetSequence(ctx, seq)
		if err != nil {
			return err
		}
	}

	return nil
}