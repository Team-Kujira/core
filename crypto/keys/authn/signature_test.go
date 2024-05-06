package authn

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"testing"

	cometcrypto "github.com/cometbft/cometbft/crypto"
	"github.com/stretchr/testify/require"
)

// CollectedClientData represents the contextual bindings of both the WebAuthn Relying Party
// and the client. It is a key-value mapping whose keys are strings. Values can be any type
// that has a valid encoding in JSON. Its structure is defined by the following Web IDL.
// https://www.w3.org/TR/webauthn/#sec-client-data
type CollectedClientData struct {
	// Type the string "webauthn.create" when creating new credentials,
	// and "webauthn.get" when getting an assertion from an existing credential. The
	// purpose of this member is to prevent certain types of signature confusion attacks
	//(where an attacker substitutes one legitimate signature for another).
	Type         CeremonyType  `json:"type"`
	Challenge    string        `json:"challenge"`
	Origin       string        `json:"origin"`
	TokenBinding *TokenBinding `json:"tokenBinding,omitempty"`
	// Chromium (Chrome) returns a hint sometimes about how to handle clientDataJSON in a safe manner
	Hint string `json:"new_keys_may_be_added_here,omitempty"`
}

type TokenBinding struct {
	Status TokenBindingStatus `json:"status"`
	ID     string             `json:"id,omitempty"`
}

type TokenBindingStatus string

type CeremonyType string

const (
	CreateCeremony CeremonyType = "webauthn.create"
	AssertCeremony CeremonyType = "webauthn.get"
)

// func (pubKey *PubKey) VerifySignature(msg []byte, sigStr []byte) bool {
// 	type CBORSignature struct {
// 		AuthenticatorData string `json:"authenticatorData"`
// 		ClientDataJSON    string `json:"clientDataJSON"`
// 		Signature         string `json:"signature"`
// 	}

// 	cborSig := CBORSignature{}
// 	err := json.Unmarshal(sigStr, &cborSig)
// 	if err != nil {
// 		return false
// 	}

// 	clientDataJSON, err := hex.DecodeString(cborSig.ClientDataJSON)
// 	if err != nil {
// 		return false
// 	}

// 	clientData := make(map[string]interface{})
// 	err = json.Unmarshal(clientDataJSON, &clientData)
// 	if err != nil {
// 		return false
// 	}

// 	challenge, err := base64.RawURLEncoding.DecodeString(clientData["challenge"].(string))
// 	if err != nil {
// 		return false
// 	}

// 	// Check challenge == msg
// 	if !bytes.Equal(challenge, msg) {
// 		return false
// 	}

// 	publicKey := &cecdsa.PublicKey{Curve: elliptic.P256()}
// 	publicKey.X, publicKey.Y = elliptic.UnmarshalCompressed(elliptic.P256(), pubKey.Key)
// 	if publicKey.X == nil || publicKey.Y == nil {
// 		return false
// 	}

// 	type ECDSASignature struct {
// 		R, S *big.Int
// 	}

// 	signatureBytes, err := hex.DecodeString(cborSig.Signature)
// 	if err != nil {
// 		return false
// 	}

// 	e := &ECDSASignature{}
// 	_, err = asn1.Unmarshal(signatureBytes, e)
// 	if err != nil {
// 		return false
// 	}

// 	authenticatorData, err := hex.DecodeString(cborSig.AuthenticatorData)
// 	if err != nil {
// 		return false
// 	}

// 	clientDataHash := sha256.Sum256(clientDataJSON)
// 	payload := append(authenticatorData, clientDataHash[:]...)

// 	h := crypto.SHA256.New()
// 	h.Write(payload)

// 	return cecdsa.Verify(publicKey, h.Sum(nil), e.R, e.S)
// }

// import (
// 	"crypto/ecdsa"
// 	"crypto/elliptic"
// 	"crypto/rand"
// 	"crypto/sha256"
// 	"fmt"
// )

// func Example() {
// 	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
// 	if err != nil {
// 		panic(err)
// 	}

// 	msg := "hello, world"
// 	hash := sha256.Sum256([]byte(msg))

// 	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("signature: %x\n", sig)

// 	valid := ecdsa.VerifyASN1(&privateKey.PublicKey, hash[:], sig)
// 	fmt.Println("signature verified:", valid)
// }

func TestVerifySignature(t *testing.T) {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	require.NoError(t, err)
	pkBytes := elliptic.MarshalCompressed(curve, privateKey.PublicKey.X, privateKey.PublicKey.Y)
	pk := PubKey{
		KeyId: "a099eda0fb05e5783379f73a06acca726673b8e07e436edcd0d71645982af65c",
		Key:   pkBytes,
	}

	authenticatorData := cometcrypto.CRandBytes(37)
	msg := cometcrypto.CRandBytes(1000)
	clientData := CollectedClientData{
		// purpose of this member is to prevent certain types of signature confusion attacks
		//(where an attacker substitutes one legitimate signature for another).
		Type:         "webauthn.create",
		Challenge:    base64.RawURLEncoding.EncodeToString(msg),
		Origin:       "https://blue.kujira.network",
		TokenBinding: nil,
		Hint:         "",
	}
	clientDataJSON, err := json.Marshal(clientData)
	require.NoError(t, err)
	clientDataHash := sha256.Sum256(clientDataJSON)
	payload := append(authenticatorData, clientDataHash[:]...)

	h := crypto.SHA256.New()
	h.Write(payload)
	digest := h.Sum(nil)

	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, digest)
	require.NoError(t, err)

	cborSig := CBORSignature{
		AuthenticatorData: hex.EncodeToString(authenticatorData),
		ClientDataJSON:    hex.EncodeToString(clientDataJSON),
		Signature:         hex.EncodeToString(sig),
	}

	sigBytes, err := json.Marshal(cborSig)
	require.NoError(t, err)
	require.True(t, pk.VerifySignature(msg, sigBytes))

	// Mutate the signature
	for i := range sig {
		sigCpy := make([]byte, len(sig))
		copy(sigCpy, sig)
		sigCpy[i] ^= byte(i + 1)
		require.False(t, pk.VerifySignature(msg, sigCpy))
	}

	// Mutate the message
	msg[1] ^= byte(2)
	require.False(t, pk.VerifySignature(msg, sig))
}
