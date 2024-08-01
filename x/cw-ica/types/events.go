package types

const (
	EventTypeICATxCallbackFailure       = "ica_tx_callback_failure"
	EventTypeICARegisterCallbackFailure = "ica_register_callback_failure"
	EventTypeICATimeoutCallbackFailure  = "ica_timeout_callback_failure"
)

// event attributes returned
const (
	AttributePacketSourcePort    = "source_port"
	AttributePacketSourceChannel = "source_channel"
	AttributePacketSequence      = "sequence"
)
