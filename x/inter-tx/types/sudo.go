package types

type MessageCallback struct {
	Callback Callback `json:"callback"`
}

type Callback struct {
	ConnId     string `json:"connection_id"`
	AccId      string `json:"account_id"`
	TxId       uint64 `json:"tx_id"`
	ResultCode uint64 `json:"result_code"` // Success(0) | Failure(1) | Timeout(2)
	ResultData []byte `json:"result_data"`
}

var (
	ResultCodeSuccess uint64 = 0
	ResultCodeFailure uint64 = 1
	ResultCodeTimeout uint64 = 2
)
