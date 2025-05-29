package inbox

import (
	"context"

	"github.com/thealiakbari/hichapp/pkg/common/kafka/inbox/domain/dto"
)

type InboxMessageService interface {
	Create(ctx context.Context, dot dto.InboxMessageCreateDto) (dto.InboxMessageGetDto, error)
	FindByIdAndAggregateType(ctx context.Context, id, aggregateType string) (*dto.InboxMessageGetDto, error)
}
