package kafka

import (
	"errors"

	"github.com/google/uuid"
	"github.com/thealiakbari/hichapp/pkg/common/db"
	"github.com/thealiakbari/hichapp/pkg/common/kafka/outbox"
	"github.com/thealiakbari/hichapp/pkg/common/kafka/outbox/domain/dto"
	"github.com/thealiakbari/hichapp/pkg/common/kafka/outbox/domain/repository"
	"github.com/thealiakbari/hichapp/pkg/common/kafka/outbox/service"
	"github.com/thealiakbari/hichapp/pkg/common/logger"
	"github.com/thealiakbari/hichapp/pkg/common/middleware"
	"golang.org/x/net/context"
)

var TraceIdError = errors.New("trace-id is not set in context")

type OutboxConfig struct {
	Db  db.DBWrapper
	Log logger.Logger
}

type OutboxWrapper interface {
	Put(ctx context.Context, aggregateId string, msg outboxMessage) (dto.OutboxMessageGetDto, error)
}

type outboxWrapperService struct {
	db            db.DBWrapper
	outboxService outbox.OutboxMessageService
	log           logger.Logger
}

func NewOutboxWrapper(config OutboxConfig) OutboxWrapper {
	instance := outboxWrapperService{
		db: config.Db,
	}

	instance.log = config.Log.ForService(outboxWrapperService{})

	err := outbox.OutboxCheckup(instance.db)
	if err != nil {
		instance.log.Panicf(nil, "Checking or migration of the OutboxMessages failed", logger.Error(err))
	}

	outboxRepo := repository.NewInboxMessageRepository(instance.db.DB)
	outboxSvc := service.NewOutboxMessageService(outboxRepo)
	instance.outboxService = outboxSvc

	return instance
}

func (o outboxWrapperService) Put(ctx context.Context, aggregateId string, msg outboxMessage) (dto.OutboxMessageGetDto, error) {
	traceIdUuid, ok := ctx.Value(middleware.TraceIdKey).(uuid.UUID)
	if !ok {
		o.log.MethodError(ctx, msg, "context has not trace-id")
		return dto.OutboxMessageGetDto{}, TraceIdError
	}

	outMsg, err := o.outboxService.Create(ctx, dto.OutboxMessageCreateDto{
		TraceId:       traceIdUuid.String(),
		AggregateId:   aggregateId,
		AggregateType: msg.topic,
		Type:          msg.messageType,
		Name:          msg.name,
		Payload:       msg.payload,
	})
	if err != nil {
		o.log.MethodError(ctx, msg, "failed to create outbox message", logger.Error(err))
		return dto.OutboxMessageGetDto{}, err
	}

	return outMsg, nil
}
