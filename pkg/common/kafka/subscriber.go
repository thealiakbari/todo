package kafka

import (
	"context"

	wk "github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
)

type subscriber struct {
	kafka      Kafka
	subscriber message.Subscriber
}

func (k Kafka) newSubscriber() (message.Subscriber, error) {
	subs, err := wk.NewSubscriber(
		wk.SubscriberConfig{
			Brokers:               k.brokers,
			ConsumerGroup:         k.groupId,
			Unmarshaler:           wk.DefaultMarshaler{},
			OverwriteSaramaConfig: k.saramaSubscriberConfig,
		},
		k.watermillLogger,
	)
	if err != nil {
		return nil, err
	}

	return subscriber{
		k,
		subs,
	}, nil
}

func (s subscriber) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	return s.subscriber.Subscribe(ctx, topic)
}

func (s subscriber) Close() error {
	return s.subscriber.Close()
}
