package service

import (
	"context"

	"github.com/thealiakbari/hichapp/pkg/common/kafka/inbox"
	inboxDto "github.com/thealiakbari/hichapp/pkg/common/kafka/inbox/domain/dto"
	"github.com/thealiakbari/hichapp/pkg/common/kafka/inbox/domain/repository"
	"github.com/thealiakbari/hichapp/pkg/common/kafka/inbox/domain/transform"
)

type inboxMessageService struct {
	inboxRepository repository.InboxMessagesRepository
}

func NewInboxMessageService(repo repository.InboxMessagesRepository) inbox.InboxMessageService {
	return inboxMessageService{
		repo,
	}
}

func (i inboxMessageService) Create(ctx context.Context, dto inboxDto.InboxMessageCreateDto) (inboxDto.InboxMessageGetDto, error) {
	err := dto.Validate()
	if err != nil {
		return inboxDto.InboxMessageGetDto{}, err
	}

	entity, err := i.inboxRepository.Create(ctx, transform.InboxCreateDtoToEntity(dto))
	if err != nil {
		return inboxDto.InboxMessageGetDto{}, err
	}
	res := transform.InboxEntityToGetDto(&entity)
	return *res, nil
}

func (i inboxMessageService) FindByIdAndAggregateType(ctx context.Context, id, aggregateType string) (*inboxDto.InboxMessageGetDto, error) {
	entity, err := i.inboxRepository.FindByIdAndAggregateType(ctx, id, aggregateType)
	if err != nil {
		return nil, err
	}

	var res *inboxDto.InboxMessageGetDto

	if entity != nil {
		res = transform.InboxEntityToGetDto(entity)
	}

	return res, nil
}
