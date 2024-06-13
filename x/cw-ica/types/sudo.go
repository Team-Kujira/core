package types

type MessageTxCallback struct {
	IcaTxCallback IcaTxCallbackData `json:"ica_tx_callback"`
}

type MessageRegisterCallback struct {
	IcaRegisterCallback IcaRegisterCallbackData `json:"ica_register_callback"`
}

type IcaRegisterCallbackData struct {
	ConnID   string            `json:"connection_id"`
	AccID    string            `json:"account_id"`
	Callback []byte            `json:"callback"`
	Result   IcaCallbackResult `json:"result"`
}

type IcaTxCallbackData struct {
	ConnID   string            `json:"connection_id"`
	AccID    string            `json:"account_id"`
	Sequence uint64            `json:"sequence"`
	Callback []byte            `json:"callback"`
	Result   IcaCallbackResult `json:"result"`
}

type IcaCallbackResult struct {
	Success *IcaCallbackSuccess `json:"success,omitempty"`
	Error   *IcaCallbackError   `json:"error,omitempty"`
	Timeout *IcaCallbackTimeout `json:"timeout,omitempty"`
}

type IcaCallbackSuccess struct {
	Data []byte `json:"data"`
}

type IcaCallbackError struct {
	Error string `json:"error"`
}

type IcaCallbackTimeout struct{}

type MessageTransferCallback struct {
	TransferCallback TransferCallbackData `json:"transfer_callback"`
}

type TransferCallbackData struct {
	Port     string            `json:"port"`
	Channel  string            `json:"channel"`
	Sequence uint64            `json:"sequence"`
	Receiver string            `json:"receiver"`
	Denom    string            `json:"denom"`
	Amount   string            `json:"amount"`
	Memo     string            `json:"memo"`
	Result   IcaCallbackResult `json:"result"`
	Callback []byte            `json:"callback"`
}

type MessageTransferReceipt struct {
	TransferReceipt TransferReceiptData `json:"transfer_receipt"`
}

type TransferReceiptData struct {
	Port     string `json:"port"`
	Channel  string `json:"channel"`
	Sequence uint64 `json:"sequence"`
	Sender   string `json:"sender"`
	Denom    string `json:"denom"`
	Amount   string `json:"amount"`
	Memo     string `json:"memo"`
}
