package codec

import (
	authn "github.com/Team-Kujira/core/crypto/keys/authn"

	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCrypto registers all crypto dependency types with the provided Amino
// codec.
func RegisterCrypto(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(authn.PubKey{},
		authn.PubKeyName, nil)
}
