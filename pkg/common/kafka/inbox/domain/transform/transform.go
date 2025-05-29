package transform

import (
	"github.com/thealiakbari/hichapp/pkg/common/kafka/inbox/domain/dto"
	"github.com/thealiakbari/hichapp/pkg/common/kafka/inbox/domain/entity"
)

func InboxCreateDtoToEntity(in dto.InboxMessageCreateDto) entity.InboxMessage {
	return entity.InboxMessage{
		Id:            in.Id,
		AggregateId:   in.AggregateId,
		Type:          in.Type,
		CorrelationId: in.CorrelationId,
		TraceId:       in.TraceId,
		Payload:       in.Payload,
		State:         in.State,
		Status:        in.Status,
		RetryCount:    in.RetryCount,
		WaitDuration:  in.WaitDuration,
		Metadata:      in.Metadata,
	}
}

func InboxEntityToGetDto(in *entity.InboxMessage) *dto.InboxMessageGetDto {
	return &dto.InboxMessageGetDto{
		Id:            in.Id,
		AggregateId:   in.AggregateId,
		Type:          in.Type,
		CorrelationId: in.CorrelationId,
		TraceId:       in.TraceId,
		Payload:       in.Payload,
		State:         in.State,
		Status:        in.Status,
		Metadata:      in.Metadata,
		RetryCount:    in.RetryCount,
		WaitDuration:  in.WaitDuration,
		CreatedAt:     in.CreatedAt,
	}
}
