package codec

import (
	"github.com/Team-Kujira/core/crypto/keys/ecdsa256"

	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCrypto registers all crypto dependency types with the provided Amino
// codec.
func RegisterCrypto(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(ecdsa256.PubKey{},
		ecdsa256.PubKeyName, nil)
}
