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

func GenerateAuthnKey(t *testing.T) (*ecdsa.PrivateKey, PubKey) {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	require.NoError(t, err)
	pkBytes := elliptic.MarshalCompressed(curve, privateKey.PublicKey.X, privateKey.PublicKey.Y)
	pk := PubKey{
		KeyId: "a099eda0fb05e5783379f73a06acca726673b8e07e436edcd0d71645982af65c",
		Key:   pkBytes,
	}
	return privateKey, pk
}

func GenerateClientData(t *testing.T, msg []byte) []byte {
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
	return clientDataJSON
}

func TestVerifySignature(t *testing.T) {
	privateKey, pk := GenerateAuthnKey(t)
	authenticatorData := cometcrypto.CRandBytes(37)
	msg := cometcrypto.CRandBytes(1000)
	clientDataJSON := GenerateClientData(t, msg)
	clientDataHash := sha256.Sum256(clientDataJSON)
	payload := append(authenticatorData, clientDataHash[:]...)

	h := crypto.SHA256.New()
	h.Write(payload)
	digest := h.Sum(nil)

	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, digest)
	require.NoError(t, err)

	cborSig := Signature{
		AuthenticatorData: hex.EncodeToString(authenticatorData),
		ClientDataJSON:    hex.EncodeToString(clientDataJSON),
		Signature:         hex.EncodeToString(sig),
	}

	sigBytes, err := json.Marshal(cborSig)
	require.NoError(t, err)
	require.True(t, pk.VerifySignature(msg, sigBytes))

	// Mutate the signature
	for i := range sigBytes {
		sigCpy := make([]byte, len(sigBytes))
		copy(sigCpy, sigBytes)
		sigCpy[i] ^= byte(0xFF)
		require.False(t, pk.VerifySignature(msg, sigCpy))
	}

	// Mutate the message
	msg[1] ^= byte(2)
	require.False(t, pk.VerifySignature(msg, sigBytes))
}

func TestVerifySignature_ChallengeStdEncoding(t *testing.T) {
	privateKey, pk := GenerateAuthnKey(t)
	authenticatorData := cometcrypto.CRandBytes(37)
	msg := cometcrypto.CRandBytes(1000)
	clientData := CollectedClientData{
		// purpose of this member is to prevent certain types of signature confusion attacks
		//(where an attacker substitutes one legitimate signature for another).
		Type:         "webauthn.create",
		Challenge:    base64.StdEncoding.EncodeToString(msg),
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

	cborSig := Signature{
		AuthenticatorData: hex.EncodeToString(authenticatorData),
		ClientDataJSON:    hex.EncodeToString(clientDataJSON),
		Signature:         hex.EncodeToString(sig),
	}

	sigBytes, err := json.Marshal(cborSig)
	require.NoError(t, err)
	require.False(t, pk.VerifySignature(msg, sigBytes))
}

func TestVerifySignature_ChallengeHexEncoding(t *testing.T) {
	privateKey, pk := GenerateAuthnKey(t)
	authenticatorData := cometcrypto.CRandBytes(37)
	msg := cometcrypto.CRandBytes(1000)
	clientData := CollectedClientData{
		// purpose of this member is to prevent certain types of signature confusion attacks
		//(where an attacker substitutes one legitimate signature for another).
		Type:         "webauthn.create",
		Challenge:    hex.EncodeToString(msg),
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

	cborSig := Signature{
		AuthenticatorData: hex.EncodeToString(authenticatorData),
		ClientDataJSON:    hex.EncodeToString(clientDataJSON),
		Signature:         hex.EncodeToString(sig),
	}

	sigBytes, err := json.Marshal(cborSig)
	require.NoError(t, err)
	require.False(t, pk.VerifySignature(msg, sigBytes))
}

func TestVerifySignature_ChallengeEmpty(t *testing.T) {
	privateKey, pk := GenerateAuthnKey(t)
	authenticatorData := cometcrypto.CRandBytes(37)
	msg := cometcrypto.CRandBytes(1000)
	clientData := CollectedClientData{
		// purpose of this member is to prevent certain types of signature confusion attacks
		//(where an attacker substitutes one legitimate signature for another).
		Type:         "webauthn.create",
		Challenge:    "",
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

	cborSig := Signature{
		AuthenticatorData: hex.EncodeToString(authenticatorData),
		ClientDataJSON:    hex.EncodeToString(clientDataJSON),
		Signature:         hex.EncodeToString(sig),
	}

	sigBytes, err := json.Marshal(cborSig)
	require.NoError(t, err)
	require.False(t, pk.VerifySignature(msg, sigBytes))
}

func TestVerifySignature_ChallengeNil(t *testing.T) {
	privateKey, pk := GenerateAuthnKey(t)
	authenticatorData := cometcrypto.CRandBytes(37)
	msg := cometcrypto.CRandBytes(1000)
	clientData := map[string]interface{}{
		// purpose of this member is to prevent certain types of signature confusion attacks
		//(where an attacker substitutes one legitimate signature for another).
		"type":   "webauthn.create",
		"origin": "https://blue.kujira.network",
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

	cborSig := Signature{
		AuthenticatorData: hex.EncodeToString(authenticatorData),
		ClientDataJSON:    hex.EncodeToString(clientDataJSON),
		Signature:         hex.EncodeToString(sig),
	}

	sigBytes, err := json.Marshal(cborSig)
	require.NoError(t, err)
	require.False(t, pk.VerifySignature(msg, sigBytes))
}

func TestVerifySignature_ChallengeInteger(t *testing.T) {
	privateKey, pk := GenerateAuthnKey(t)
	authenticatorData := cometcrypto.CRandBytes(37)
	msg := cometcrypto.CRandBytes(1000)
	clientData := map[string]interface{}{
		// purpose of this member is to prevent certain types of signature confusion attacks
		//(where an attacker substitutes one legitimate signature for another).
		"type":      "webauthn.create",
		"origin":    "https://blue.kujira.network",
		"challenge": 1,
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

	cborSig := Signature{
		AuthenticatorData: hex.EncodeToString(authenticatorData),
		ClientDataJSON:    hex.EncodeToString(clientDataJSON),
		Signature:         hex.EncodeToString(sig),
	}

	sigBytes, err := json.Marshal(cborSig)
	require.NoError(t, err)
	require.False(t, pk.VerifySignature(msg, sigBytes))
}

func TestVerifySignature_ClientDataJSONEmpty(t *testing.T) {
	privateKey, pk := GenerateAuthnKey(t)
	authenticatorData := cometcrypto.CRandBytes(37)
	msg := cometcrypto.CRandBytes(1000)
	payload := authenticatorData

	h := crypto.SHA256.New()
	h.Write(payload)
	digest := h.Sum(nil)

	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, digest)
	require.NoError(t, err)

	type CBORSignature struct {
		AuthenticatorData string `json:"authenticatorData"`
		Signature         string `json:"signature"`
	}
	cborSig := CBORSignature{
		AuthenticatorData: hex.EncodeToString(authenticatorData),
		Signature:         hex.EncodeToString(sig),
	}

	sigBytes, err := json.Marshal(cborSig)
	require.NoError(t, err)
	require.False(t, pk.VerifySignature(msg, sigBytes))
}

func TestVerifySignature_UncompressedPubKey(t *testing.T) {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	require.NoError(t, err)
	pkBytes := elliptic.Marshal(curve, privateKey.PublicKey.X, privateKey.PublicKey.Y)
	pk := PubKey{
		KeyId: "a099eda0fb05e5783379f73a06acca726673b8e07e436edcd0d71645982af65c",
		Key:   pkBytes,
	}

	authenticatorData := cometcrypto.CRandBytes(37)
	msg := cometcrypto.CRandBytes(1000)

	clientDataJSON := GenerateClientData(t, msg)
	clientDataHash := sha256.Sum256(clientDataJSON)
	payload := append(authenticatorData, clientDataHash[:]...)

	h := crypto.SHA256.New()
	h.Write(payload)
	digest := h.Sum(nil)

	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, digest)
	require.NoError(t, err)

	cborSig := Signature{
		AuthenticatorData: hex.EncodeToString(authenticatorData),
		ClientDataJSON:    hex.EncodeToString(clientDataJSON),
		Signature:         hex.EncodeToString(sig),
	}

	sigBytes, err := json.Marshal(cborSig)
	require.NoError(t, err)
	require.False(t, pk.VerifySignature(msg, sigBytes))
}

func TestVerifySignature_AnotherPubKey(t *testing.T) {
	privateKey, _ := GenerateAuthnKey(t)
	_, pk := GenerateAuthnKey(t)

	authenticatorData := cometcrypto.CRandBytes(37)
	msg := cometcrypto.CRandBytes(1000)

	clientDataJSON := GenerateClientData(t, msg)
	clientDataHash := sha256.Sum256(clientDataJSON)
	payload := append(authenticatorData, clientDataHash[:]...)

	h := crypto.SHA256.New()
	h.Write(payload)
	digest := h.Sum(nil)

	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, digest)
	require.NoError(t, err)

	cborSig := Signature{
		AuthenticatorData: hex.EncodeToString(authenticatorData),
		ClientDataJSON:    hex.EncodeToString(clientDataJSON),
		Signature:         hex.EncodeToString(sig),
	}

	sigBytes, err := json.Marshal(cborSig)
	require.NoError(t, err)
	require.False(t, pk.VerifySignature(msg, sigBytes))
}

func TestVerifySignature_SignaureNotInASN1(t *testing.T) {
	privateKey, pk := GenerateAuthnKey(t)
	authenticatorData := cometcrypto.CRandBytes(37)
	msg := cometcrypto.CRandBytes(1000)

	clientDataJSON := GenerateClientData(t, msg)
	clientDataHash := sha256.Sum256(clientDataJSON)
	payload := append(authenticatorData, clientDataHash[:]...)

	h := crypto.SHA256.New()
	h.Write(payload)
	digest := h.Sum(nil)

	sigR, sigS, err := ecdsa.Sign(rand.Reader, privateKey, digest)
	require.NoError(t, err)
	sig := append(sigR.Bytes(), sigS.Bytes()...)

	cborSig := Signature{
		AuthenticatorData: hex.EncodeToString(authenticatorData),
		ClientDataJSON:    hex.EncodeToString(clientDataJSON),
		Signature:         hex.EncodeToString(sig),
	}

	sigBytes, err := json.Marshal(cborSig)
	require.NoError(t, err)
	require.False(t, pk.VerifySignature(msg, sigBytes))
}

func TestVerifySignature_EmptyAuthenticatorData(t *testing.T) {
	privateKey, pk := GenerateAuthnKey(t)
	authenticatorData := cometcrypto.CRandBytes(37)
	msg := cometcrypto.CRandBytes(1000)

	clientDataJSON := GenerateClientData(t, msg)
	clientDataHash := sha256.Sum256(clientDataJSON)
	payload := append(authenticatorData, clientDataHash[:]...)

	h := crypto.SHA256.New()
	h.Write(payload)
	digest := h.Sum(nil)

	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, digest)
	require.NoError(t, err)

	cborSig := Signature{
		AuthenticatorData: "",
		ClientDataJSON:    hex.EncodeToString(clientDataJSON),
		Signature:         hex.EncodeToString(sig),
	}

	sigBytes, err := json.Marshal(cborSig)
	require.NoError(t, err)
	require.False(t, pk.VerifySignature(msg, sigBytes))
}

func TestVerifySignature_SignatureEncodingInBase64(t *testing.T) {
	privateKey, pk := GenerateAuthnKey(t)
	authenticatorData := cometcrypto.CRandBytes(37)
	msg := cometcrypto.CRandBytes(1000)

	clientDataJSON := GenerateClientData(t, msg)
	clientDataHash := sha256.Sum256(clientDataJSON)
	payload := append(authenticatorData, clientDataHash[:]...)

	h := crypto.SHA256.New()
	h.Write(payload)
	digest := h.Sum(nil)

	sigR, sigS, err := ecdsa.Sign(rand.Reader, privateKey, digest)
	require.NoError(t, err)
	sig := append(sigR.Bytes(), sigS.Bytes()...)

	cborSig := Signature{
		AuthenticatorData: hex.EncodeToString(authenticatorData),
		ClientDataJSON:    hex.EncodeToString(clientDataJSON),
		Signature:         base64.StdEncoding.EncodeToString(sig),
	}

	sigBytes, err := json.Marshal(cborSig)
	require.NoError(t, err)
	require.False(t, pk.VerifySignature(msg, sigBytes))
}

func TestVerifySignature_ClientDataJSONEncodingInBase64(t *testing.T) {
	privateKey, pk := GenerateAuthnKey(t)
	authenticatorData := cometcrypto.CRandBytes(37)
	msg := cometcrypto.CRandBytes(1000)
	clientDataJSON := GenerateClientData(t, msg)
	clientDataHash := sha256.Sum256(clientDataJSON)
	payload := append(authenticatorData, clientDataHash[:]...)

	h := crypto.SHA256.New()
	h.Write(payload)
	digest := h.Sum(nil)

	sigR, sigS, err := ecdsa.Sign(rand.Reader, privateKey, digest)
	require.NoError(t, err)
	sig := append(sigR.Bytes(), sigS.Bytes()...)

	cborSig := Signature{
		AuthenticatorData: hex.EncodeToString(authenticatorData),
		ClientDataJSON:    base64.StdEncoding.EncodeToString(clientDataJSON),
		Signature:         hex.EncodeToString(sig),
	}

	sigBytes, err := json.Marshal(cborSig)
	require.NoError(t, err)
	require.False(t, pk.VerifySignature(msg, sigBytes))
}

func TestVerifySignature_EmptySignature(t *testing.T) {
	_, pk := GenerateAuthnKey(t)
	msg := cometcrypto.CRandBytes(1000)
	require.False(t, pk.VerifySignature(msg, []byte{}))
}

func TestVerifySignature_AuthenticatorDataLength(t *testing.T) {
	privateKey, pk := GenerateAuthnKey(t)
	authenticatorData := cometcrypto.CRandBytes(10)
	msg := cometcrypto.CRandBytes(1000)
	clientDataJSON := GenerateClientData(t, msg)
	clientDataHash := sha256.Sum256(clientDataJSON)
	payload := append(authenticatorData, clientDataHash[:]...)

	h := crypto.SHA256.New()
	h.Write(payload)
	digest := h.Sum(nil)

	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, digest)
	require.NoError(t, err)

	cborSig := Signature{
		AuthenticatorData: hex.EncodeToString(authenticatorData),
		ClientDataJSON:    hex.EncodeToString(clientDataJSON),
		Signature:         hex.EncodeToString(sig),
	}

	sigBytes, err := json.Marshal(cborSig)
	require.NoError(t, err)
	require.False(t, pk.VerifySignature(msg, sigBytes))
}
