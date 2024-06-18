package authn

import (
	// "bytes"
	"bytes"
	"crypto"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"math/big"

	ecdsa "crypto/ecdsa"
)

type CBORSignature struct {
	AuthenticatorData string `json:"authenticatorData"`
	ClientDataJSON    string `json:"clientDataJSON"`
	Signature         string `json:"signature"`
}

// VerifyBytes verifies a signature of the form R || S.
// It rejects signatures which are not in lower-S form.
// See https://github.com/Team-Kujira/kujira.js/blob/master/src/authn/AuthnWebSigner.ts for a reference implementation
// signing a transaction with the webauthn API
func (pubKey *PubKey) VerifySignature(msg []byte, sigStr []byte) bool {
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

	challengeBase64, ok := clientData["challenge"].(string)
	if !ok {
		return false
	}
	challenge, err := base64.RawURLEncoding.DecodeString(challengeBase64)
	if err != nil {
		return false
	}

	// Check challenge == msg
	if !bytes.Equal(challenge, msg) {
		return false
	}

	publicKey := &ecdsa.PublicKey{Curve: elliptic.P256()}
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

	return ecdsa.Verify(publicKey, h.Sum(nil), e.R, e.S)
}
