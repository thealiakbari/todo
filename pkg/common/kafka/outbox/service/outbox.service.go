package service

import (
	"context"

	"github.com/thealiakbari/hichapp/pkg/common/kafka/outbox"
	outboxDto "github.com/thealiakbari/hichapp/pkg/common/kafka/outbox/domain/dto"
	"github.com/thealiakbari/hichapp/pkg/common/kafka/outbox/domain/repository"
	"github.com/thealiakbari/hichapp/pkg/common/kafka/outbox/domain/transform"
)

type outboxMessageService struct {
	outboxRepository repository.OutboxMessagesRepository
}

func NewOutboxMessageService(repo repository.OutboxMessagesRepository) outbox.OutboxMessageService {
	return outboxMessageService{
		repo,
	}
}

func (i outboxMessageService) Create(ctx context.Context, dto outboxDto.OutboxMessageCreateDto) (outboxDto.OutboxMessageGetDto, error) {
	err := dto.Validate()
	if err != nil {
		return outboxDto.OutboxMessageGetDto{}, err
	}

	entity, err := i.outboxRepository.Create(ctx, transform.OutboxCreateDtoToEntity(dto))
	if err != nil {
		return outboxDto.OutboxMessageGetDto{}, err
	}
	res := transform.OutboxEntityToGetDto(&entity)
	return *res, nil
}
