package ecdsa256

import (
	// "bytes"
	"bytes"
	"crypto"

	cecdsa "crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"math/big"

	secp256k1 "github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"

	cometcrypto "github.com/cometbft/cometbft/crypto"
)

// Sign creates an ECDSA signature on curve Secp256k1, using SHA256 on the msg.
// The returned signature will be of the form R || S (in lower-S form).
func (privKey *PrivKey) Sign(msg []byte) ([]byte, error) {
	priv, _ := secp256k1.PrivKeyFromBytes(privKey.Key)
	sig, err := ecdsa.SignCompact(priv, cometcrypto.Sha256(msg), false)
	if err != nil {
		return nil, err
	}
	// remove the first byte which is compactSigRecoveryCode
	return sig[1:], nil
}

// VerifyBytes verifies a signature of the form R || S.
// It rejects signatures which are not in lower-S form.
func (pubKey *PubKey) VerifySignature(msg []byte, sigStr []byte) bool {
	// return false
	type CBORSignature struct {
		AuthenticatorData string `json:"authenticatorData"`
		ClientDataJSON    string `json:"clientDataJSON"`
		Signature         string `json:"signature"`
	}

	cborSig := CBORSignature{}
	err := json.Unmarshal(sigStr, &cborSig)
	if err != nil {
		return false
	}

	clientDataJSON, err := hex.DecodeString(cborSig.ClientDataJSON)
	if err != nil {
		return false
	}

	clientData := make(map[string]interface{})
	err = json.Unmarshal(clientDataJSON, &clientData)
	if err != nil {
		return false
	}

	challenge, err := base64.RawURLEncoding.DecodeString(clientData["challenge"].(string))
	if err != nil {
		return false
	}

	// Check challenge == msg
	if !bytes.Equal(challenge, msg) {
		return false
	}

	publicKey := &cecdsa.PublicKey{Curve: elliptic.P256()}
	publicKey.X, publicKey.Y = elliptic.UnmarshalCompressed(elliptic.P256(), pubKey.Key)
	if publicKey.X == nil || publicKey.Y == nil {
		return false
	}

	type ECDSASignature struct {
		R, S *big.Int
	}

	signatureBytes, err := hex.DecodeString(cborSig.Signature)
	if err != nil {
		return false
	}

	e := &ECDSASignature{}
	_, err = asn1.Unmarshal(signatureBytes, e)
	if err != nil {
		return false
	}

	authenticatorData, err := hex.DecodeString(cborSig.AuthenticatorData)
	if err != nil {
		return false
	}

	clientDataHash := sha256.Sum256(clientDataJSON)
	payload := append(authenticatorData, clientDataHash[:]...)

	h := crypto.SHA256.New()
	h.Write(payload)

	return cecdsa.Verify(publicKey, h.Sum(nil), e.R, e.S)
}
