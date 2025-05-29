package entity

import (
	"time"

	"github.com/thealiakbari/hichapp/pkg/common/kafka/inbox/domain/enum"
)

type InboxMessage struct {
	Id            string           `gorm:"column:id;type:uuid"`
	AggregateId   string           `gorm:"column:aggregate_id;type:uuid;not null"`
	Type          string           `gorm:"column:aggregate_type;type:varchar(100)"`
	CorrelationId string           `gorm:"column:correlation_id;type:varchar(40);not null"`
	TraceId       string           `gorm:"column:trace_id;type:varchar(40)"`
	Payload       *string          `gorm:"column:payload;type:jsonb;not null"`
	State         enum.InboxState  `gorm:"column:state;type:varchar(100)"`
	Status        enum.InboxStatus `gorm:"column:status;type:varchar(100)"`
	RetryCount    int              `gorm:"column:retry_count;type:int;not null"`
	WaitDuration  *int             `gorm:"column:wait_duration;type:int"`
	Metadata      *string          `gorm:"column:metadata;type:jsonb;not null"`
	Version       int              `gorm:"column:version;type:int;not null"`
	CreatedAt     time.Time        `gorm:"column:created_at;not null;index"`
}
