package app

import (
	"github.com/Team-Kujira/core/app/params"
	kujiracryptocodec "github.com/Team-Kujira/core/crypto/codec"
	"github.com/cosmos/cosmos-sdk/std"
)

// MakeEncodingConfig creates a new EncodingConfig with all modules registered
func MakeEncodingConfig() params.EncodingConfig {
	encodingConfig := params.MakeEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	kujiracryptocodec.RegisterCrypto(encodingConfig.Amino)
	kujiracryptocodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}
