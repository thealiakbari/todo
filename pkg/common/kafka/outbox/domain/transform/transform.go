package transform

import (
	"github.com/thealiakbari/hichapp/pkg/common/kafka/outbox/domain/dto"
	"github.com/thealiakbari/hichapp/pkg/common/kafka/outbox/domain/entity"
)

func OutboxCreateDtoToEntity(in dto.OutboxMessageCreateDto) entity.OutboxMessage {
	return entity.OutboxMessage{
		AggregateId:   in.AggregateId,
		AggregateType: in.AggregateType,
		TraceId:       in.TraceId,
		Type:          in.Type,
		Name:          in.Name,
		Payload:       in.Payload,
	}
}

func OutboxEntityToGetDto(in *entity.OutboxMessage) *dto.OutboxMessageGetDto {
	return &dto.OutboxMessageGetDto{
		Id:            in.Id,
		AggregateId:   in.AggregateId,
		AggregateType: in.AggregateType,
		TraceId:       in.TraceId,
		Name:          in.Name,
		Type:          in.Type,
		Payload:       in.Payload,
		CreatedAt:     in.CreatedAt,
	}
}
