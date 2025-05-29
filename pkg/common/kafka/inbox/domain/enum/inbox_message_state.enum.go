package enum

type InboxState string

const (
	InboxStateNone       InboxState = "none"
	InboxStateInProgress InboxState = "in_progress"
	InboxStateCompleted  InboxState = "completed"
)
