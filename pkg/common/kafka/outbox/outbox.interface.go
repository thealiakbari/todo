package outbox

import (
	"context"

	"github.com/thealiakbari/hichapp/pkg/common/kafka/outbox/domain/dto"
)

type OutboxMessageService interface {
	Create(ctx context.Context, dot dto.OutboxMessageCreateDto) (dto.OutboxMessageGetDto, error)
}
