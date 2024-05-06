package codec

import (
	authn "github.com/Team-Kujira/core/crypto/keys/authn"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

// RegisterInterfaces registers the sdk.Tx interface.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	var pk *cryptotypes.PubKey
	registry.RegisterImplementations(pk, &authn.PubKey{})
}
