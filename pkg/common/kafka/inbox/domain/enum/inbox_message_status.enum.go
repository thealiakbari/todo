package enum

type InboxStatus string

const (
	InboxStatusNone       InboxStatus = "none"
	InboxStatusProcessing InboxStatus = "processing"
	InboxStatusWaiting    InboxStatus = "waiting"
	InboxStatusRetrying   InboxStatus = "retrying"
	InboxStatusSucceeded  InboxStatus = "succeeded"
	InboxStatusFailed     InboxStatus = "failed"
)
