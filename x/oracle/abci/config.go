package abci

import (
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cast"
)

type OracleConfig struct {
	Endpoint string `mapstructure:"endpoint"`
}

const (
	flagOracleEndpoint = "oracle.endpoint"
)

// ReadOracleConfig reads the wasm specifig configuration
func ReadOracleConfig(opts servertypes.AppOptions) (OracleConfig, error) {
	cfg := OracleConfig{}
	var err error
	if v := opts.Get(flagOracleEndpoint); v != nil {
		if cfg.Endpoint, err = cast.ToStringE(v); err != nil {
			return cfg, err
		}
	}

	return cfg, nil
}
