package app

import (
	"encoding/base64"

	abci "github.com/cometbft/cometbft/abci/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
)

func (app *App) PrepareProposal(req abci.RequestPrepareProposal) abci.ResponsePrepareProposal {
	newTxs := [][]byte{}
	for _, rawTx := range req.Txs {
		tx, err := app.txConfig.TxDecoder()(rawTx)
		if err != nil {
			continue
		}
		msgs := tx.GetMsgs()
		for _, msg := range msgs {
			switch msg := msg.(type) {
			case *banktypes.MsgSend:
			case *transfertypes.MsgTransfer:
				if msg.Memo != "" {
					newRawTx, err := base64.StdEncoding.DecodeString(msg.Memo)
					if err != nil {
						continue
					}
					_, err = app.txConfig.TxDecoder()(newRawTx)
					if err != nil {
						continue
					}
					newTxs = append(newTxs, newRawTx)
				}
			case *channeltypes.MsgRecvPacket:
				packet := msg.Packet
				var data transfertypes.FungibleTokenPacketData
				if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
					continue
				}
				if data.Memo != "" {
					newRawTx, err := base64.StdEncoding.DecodeString(data.Memo)
					if err != nil {
						continue
					}
					_, err = app.txConfig.TxDecoder()(newRawTx)
					if err != nil {
						continue
					}
					newTxs = append(newTxs, newRawTx)
				}
			}
		}
	}
	txs := append(req.Txs, newTxs...)
	return abci.ResponsePrepareProposal{
		Txs: txs,
	}
}