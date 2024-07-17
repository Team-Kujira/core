package keeper

import (
	"encoding/base64"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) HandleTransferHook(ctx sdk.Context, memo string, txEncodingConfig client.TxEncodingConfig) {
	newRawTx, err := base64.StdEncoding.DecodeString(memo)
	if err != nil {
		return
	}

	tx, err := txEncodingConfig.TxDecoder()(newRawTx)
	if err != nil {
		return
	}

	cacheCtx, write := ctx.CacheContext()
	err = k.ExecuteAnte(cacheCtx, tx)
	if err != nil {
		return
	}

	_, err = k.ExecuteTxMsgs(cacheCtx, tx)
	if err == nil {
		write()
	}
}
