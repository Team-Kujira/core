package types

type MessageCallback struct {
	IcaCallback IcaCallbackData `json:"ica_callback"`
}

type IcaCallbackData struct {
	ConnId string            `json:"connection_id"`
	AccId  string            `json:"account_id"`
	TxId   uint64            `json:"tx_id"`
	Result IcaCallbackResult `json:"result"`
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

type IcaCallbackTimeout struct {
}
