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
