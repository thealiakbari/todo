package dto

import (
	"time"

	"github.com/thealiakbari/hichapp/pkg/common/kafka/outbox/domain/enum"
)

type OutboxMessageGetDto struct {
	Id            string                 `json:"id"`
	AggregateId   string                 `json:"aggregate_id"`
	AggregateType string                 `json:"aggregate_type"`
	TraceId       string                 `json:"trace_id"`
	Type          enum.OutboxMessageType `json:"type"`
	Name          string                 `json:"name"`
	Payload       *string                `json:"payload"`
	CreatedAt     time.Time              `json:"created_at"`
}
