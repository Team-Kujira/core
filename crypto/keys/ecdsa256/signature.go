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
	"fmt"
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
	fmt.Println("msg.HEX", hex.EncodeToString(msg))
	fmt.Println("VerifySignature", string(sigStr[:]))
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

	fmt.Println("CBORSignature", cborSig)

	// "clientDataJParsed": {
	//   "type": "webauthn.get",
	//   "challenge": "CpMBCpABChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEnAKLWt1amlyYTF2NDJzZnBjY2RhbnUwdmZ3eDB1eWt0NzVrYXY3eHR0dXZhejlwcRIta3VqaXJhMWsydHVsYXU3d3o2ZGE2aGZseXZtNWF1YXE4NGQ1eXB3bHQ0eGdnGhAKBXN0YWtlEgcxMDAwMDAwEnYKbgpkCh4va3VqaXJhLmNyeXB0by5lY2RzYTI1Ni5QdWJLZXkSQgpAFgGQ-0-aPaRoZOR1E712MLHeCytGYb16r70oS5x59-IlR_DXS_Emq5xotXVL28STYOVShgjII6P5o4dF0Y6MzBIECgIIARgBEgQQwJoMGghrdWppcmEtMSAC",
	//   "origin": "https://blue.kujira.network",
	//   "crossOrigin": false
	// }

	clientDataJSON, err := hex.DecodeString(cborSig.ClientDataJSON)
	if err != nil {
		return false
	}

	clientData := make(map[string]interface{})
	err = json.Unmarshal(clientDataJSON, &clientData)
	fmt.Println("clientDataErr", err)
	if err != nil {
		return false
	}
	fmt.Println("clientData", clientData)

	challenge, err := base64.RawURLEncoding.DecodeString(clientData["challenge"].(string))
	fmt.Println("clientData['challenge'].(string)", clientData["challenge"].(string))
	fmt.Println("err", err)
	if err != nil {
		return false
	}
	fmt.Println("challenge", hex.EncodeToString(challenge))

	// Check challenge == msg
	if !bytes.Equal(challenge, msg) {
		return false
	}

	// publicKey := &cecdsa.PublicKey{
	// 	Curve: elliptic.P256(),
	// 	X:     big.NewInt(0).SetBytes(pubKey.Key[:32]),
	// 	Y:     big.NewInt(0).SetBytes(pubKey.Key[32:64]),
	// }

	publicKey := &cecdsa.PublicKey{Curve: elliptic.P256()}
	publicKey.X, publicKey.Y = elliptic.UnmarshalCompressed(elliptic.P256(), pubKey.Key)
	if publicKey.X == nil || publicKey.Y == nil {
		fmt.Println("Wrong pubkey bytes", publicKey.X, publicKey.Y)
		return false
	}
	fmt.Println("X", hex.EncodeToString(publicKey.X.Bytes()))
	fmt.Println("Y", hex.EncodeToString(publicKey.Y.Bytes()))

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
	fmt.Println("clientDataHash", hex.EncodeToString(clientDataHash[:]))

	payload := append(authenticatorData, clientDataHash[:]...)
	fmt.Println("payload", hex.EncodeToString(payload[:]))

	h := crypto.SHA256.New()
	h.Write(payload)

	return cecdsa.Verify(publicKey, h.Sum(nil), e.R, e.S)
}
