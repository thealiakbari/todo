package entity

import (
	"time"

	"github.com/thealiakbari/hichapp/pkg/common/kafka/outbox/domain/enum"
)

type OutboxMessage struct {
	Id            string                 `gorm:"column:id;type:uuid;default:uuid_generate_v4()"`
	AggregateId   string                 `gorm:"column:aggregate_id;type:uuid;not null"`
	AggregateType string                 `gorm:"column:aggregate_type;type:varchar(100);not null"`
	TraceId       string                 `gorm:"column:trace_id;type:uuid;not null"`
	Type          enum.OutboxMessageType `gorm:"column:type;type:varchar(100)"`
	Name          string                 `gorm:"column:name;type:varchar(100)"`
	Payload       *string                `gorm:"column:payload;type:jsonb;not null"`
	CreatedAt     time.Time              `gorm:"column:created_at;not null;index"`
}
