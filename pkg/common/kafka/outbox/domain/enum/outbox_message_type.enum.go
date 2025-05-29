package enum

type OutboxMessageType string

const (
	OutboxMessageTypeNone     OutboxMessageType = "none"
	OutboxMessageTypeCommand  OutboxMessageType = "command"
	OutboxMessageTypeEvent    OutboxMessageType = "event"
	OutboxMessageTypeDocument OutboxMessageType = "document"
)
