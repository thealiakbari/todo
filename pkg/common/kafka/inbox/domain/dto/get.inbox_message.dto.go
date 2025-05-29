package dto

import (
	"time"

	"github.com/thealiakbari/hichapp/pkg/common/kafka/inbox/domain/enum"
)

type InboxMessageGetDto struct {
	Id            string           `json:"id"`
	AggregateId   string           `json:"aggregateId"`
	Type          string           `json:"type"`
	CorrelationId string           `json:"correlation_id"`
	TraceId       string           `json:"trace_id"`
	Payload       *string          `json:"payload"`
	State         enum.InboxState  `json:"state"`
	Status        enum.InboxStatus `json:"status"`
	Metadata      *string          `json:"metadata"`
	RetryCount    int              `json:"retry_count"`
	WaitDuration  *int             `json:"wait_duration"`
	CreatedAt     time.Time        `json:"createdAt"`
}
