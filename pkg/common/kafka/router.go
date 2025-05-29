package kafka

import (
	"context"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/thealiakbari/hichapp/pkg/common/db"
	"github.com/thealiakbari/hichapp/pkg/common/kafka/inbox/domain/dto"
	"github.com/thealiakbari/hichapp/pkg/common/kafka/inbox/domain/enum"
	"github.com/thealiakbari/hichapp/pkg/common/middleware"
	"github.com/thealiakbari/hichapp/pkg/common/utiles"
)

const defaultRouterCloseTime = 30 * time.Second

var ErrInboxTimeout = errors.New("timeout error")

const InboxHeadersKey = "headers"

type info struct {
	id            string
	aggregateId   string
	traceId       string
	correlationId string
	topic         string
}

func delayGenerator(baseDelay, maxRetries int64) func() (delay time.Duration, done bool) {
	attempt := int64(1)
	return func() (time.Duration, bool) {
		defer func() {
			attempt += 1
		}()
		return time.Duration(baseDelay * attempt), attempt >= maxRetries
	}
}

func (k Kafka) insertInbox(
	info info,
	attempt int,
	payload *string,
	processErr error,
	completed bool,
) (*dto.InboxMessageGetDto, error) {
	tx, inboxCtx, err := db.BeginTx(context.Background(), k.db.DB)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			if err := tx.Rollback().Error; err != nil {
				k.logger.Errorf("Can't inbox tx rollback:%s", err)
			}
		}
	}()

	dto := dto.InboxMessageCreateDto{
		Id:            info.id,
		AggregateId:   info.aggregateId,
		TraceId:       info.traceId,
		Type:          info.topic,
		CorrelationId: info.correlationId,
	}

	var meta *string = nil
	if processErr != nil {
		meta = utiles.Ptr(utiles.Pretty(map[string]interface{}{
			"error": processErr,
		}))
	}

	dto.RetryCount = attempt

	if payload != nil {
		dto.State = enum.InboxStateInProgress
		dto.Status = enum.InboxStatusProcessing
		dto.Payload = payload
	} else if !completed && processErr != nil {
		dto.State = enum.InboxStateInProgress
		dto.Status = enum.InboxStatusRetrying
		dto.Metadata = meta
	} else if completed && processErr == nil {
		dto.State = enum.InboxStateCompleted
		dto.Status = enum.InboxStatusSucceeded
	} else if completed && processErr != nil {
		dto.State = enum.InboxStateCompleted
		dto.Status = enum.InboxStatusFailed
		dto.Metadata = meta
	}

	inserted, err := k.inbox.Create(inboxCtx, dto)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit().Error; err != nil {
		return nil, err
	}

	return &inserted, nil
}

func (k Kafka) handleMessage(topic string, handlerFn HandlerFn, msg *message.Message) error {
	headers, err := dto.LoadInboxMessageHeader(msg.Metadata)
	if err != nil {
		meta := utiles.Pretty(msg.Metadata)
		return errors.Wrap(err, "Header validation failed, \n"+meta)
	}

	traceId, _ := uuid.Parse(headers.TraceId)
	msgCtx := context.WithValue(msg.Context(), InboxHeadersKey, *headers)
	msgCtx = context.WithValue(msgCtx, middleware.TraceIdKey, traceId)

	attempt := 1

	info := info{
		id:            headers.OutboxId,
		aggregateId:   headers.CorrelationId,
		traceId:       headers.TraceId,
		topic:         topic,
		correlationId: headers.CorrelationId,
	}

	found, err := k.inbox.FindByIdAndAggregateType(context.Background(), headers.OutboxId, topic)
	if err != nil {
		return err
	}

	if found != nil && found.Id != "" {
		k.watermillLogger.Info("Already processed InboxMessage", map[string]interface{}{
			"type":      topic,
			"lastState": utiles.Pretty(found),
		})
		return nil
	}

	payload := string(msg.Payload)

	_, err = k.insertInbox(info, attempt, &payload, nil, false)
	if err != nil {
		return err
	}

	delayGen := delayGenerator(int64(k.inboxRetry.BaseDelay), int64(k.inboxRetry.MaxRetries))
	for {
		processErr := handlerFn(msgCtx, msg.Payload)
		if processErr != nil &&
			// First attempt is not counts as a retry!
			attempt < k.inboxRetry.MaxRetries+1 {
			if !errors.Is(processErr, ErrInboxTimeout) {
				_, err := k.insertInbox(info, attempt, nil, processErr, true)
				if err != nil {
					err = errors.Wrapf(err, "Failed to insert failed process of InboxMessage: %s", processErr.Error())
					return err
				}
				return nil // send ACK
			}
			attempt += 1
			_, err := k.insertInbox(info, attempt, nil, processErr, false)
			if err != nil {
				err = errors.Wrapf(err, "Failed to insert failed and retrying process of InboxMessage: %s", processErr.Error())
				return err
			}
			delay, _ := delayGen()
			delay *= time.Duration(k.inboxRetry.ScaleFactor)
			if delay > time.Duration(k.inboxRetry.MaxDelay) {
				delay = time.Duration(k.inboxRetry.MaxDelay)
			}
			time.Sleep(delay * time.Millisecond)
			continue // another attempt
		} else if processErr != nil &&
			// First attempt is not counts as a retry!
			attempt == k.inboxRetry.MaxRetries+1 {
			_, err := k.insertInbox(info, attempt, nil, processErr, true)
			if err != nil {
				err = errors.Wrapf(err, "Failed to insert failed process of InboxMessage: %s", processErr.Error())
				return err
			}
			return nil // send ACK
		} else if processErr == nil { // simple condition just for clarification
			break
		}
	}
	_, err = k.insertInbox(info, attempt, nil, nil, true)
	if err != nil {
		return err
	}
	return nil // send ACK
}

func (k Kafka) AddHandler(handlerName, topic string, handlerFn HandlerFn) {
	routeHandler := k.router.AddNoPublisherHandler(
		handlerName,
		topic,
		k.subscriber,
		func(msg *message.Message) error {
			return k.handleMessage(topic, handlerFn, msg)
		},
	)

	if k.debugLog {
		routeHandler.AddMiddleware(func(h message.HandlerFunc) message.HandlerFunc {
			return func(message *message.Message) ([]*message.Message, error) {
				k.watermillLogger.Info(
					"Message received",
					map[string]interface{}{
						"topic":   topic,
						"headers": message.Metadata,
					},
				)

				return h(message)
			}
		})
	}
}

func (k Kafka) Run(ctx context.Context) {
	if err := k.router.Run(ctx); err != nil {
		panic(err)
	}
}

func (k Kafka) Close() error {
	return k.router.Close()
}
