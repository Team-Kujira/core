package codec

import (
	"github.com/Team-Kujira/core/crypto/keys/ecdsa256"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

// RegisterInterfaces registers the sdk.Tx interface.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	var pk *cryptotypes.PubKey
	registry.RegisterImplementations(pk, &ecdsa256.PubKey{})

	var priv *cryptotypes.PrivKey
	registry.RegisterImplementations(priv, &ecdsa256.PrivKey{})
}
