package types

import (
	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	cosmostypes "github.com/cosmos/cosmos-sdk/codec/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
)

func SerializeCosmosTx(cdc codec.BinaryCodec, msgs []*cosmostypes.Any) (bz []byte, err error) {
	// only ProtoCodec is supported
	if _, ok := cdc.(*codec.ProtoCodec); !ok {
		return nil, errors.Wrap(err, "only ProtoCodec is supported on host chain")
	}

	cosmosTx := &icatypes.CosmosTx{
		Messages: msgs,
	}

	bz, err = cdc.Marshal(cosmosTx)
	if err != nil {
		return nil, err
	}

	return bz, nil
}
