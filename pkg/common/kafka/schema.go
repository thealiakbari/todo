package kafka

import (
	"encoding/json"

	"github.com/thealiakbari/hichapp/pkg/common/kafka/outbox/domain/enum"
)

type OM interface {
	ExtractMessage()
}

type OutboxMessageSchema[T any] struct {
	Name    string
	Topic   string
	Type    enum.OutboxMessageType
	payload T
}

type outboxMessage struct {
	name        string
	topic       string
	messageType enum.OutboxMessageType
	payload     *string
}

func ToOutboxMessage[T any](schema OutboxMessageSchema[T], payload T) (outboxMessage, error) {
	bs, err := json.Marshal(payload)
	if err != nil {
		return outboxMessage{}, err
	}

	payloadStr := string(bs)
	return outboxMessage{
		name:        schema.Name,
		topic:       schema.Topic,
		payload:     &payloadStr,
		messageType: schema.Type,
	}, nil
}
