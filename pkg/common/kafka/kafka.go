package kafka

import (
	"time"

	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/thealiakbari/hichapp/pkg/common/config"
	"github.com/thealiakbari/hichapp/pkg/common/db"
	"github.com/thealiakbari/hichapp/pkg/common/kafka/inbox"
	"github.com/thealiakbari/hichapp/pkg/common/kafka/inbox/domain/repository"
	"github.com/thealiakbari/hichapp/pkg/common/kafka/inbox/service"
	"github.com/thealiakbari/hichapp/pkg/common/logger"
	"github.com/thealiakbari/hichapp/pkg/common/utiles"
)

type Kafka struct {
	logger                 logger.InfraLogger
	watermillLogger        watermill.LoggerAdapter
	brokers                []string
	saramaSubscriberConfig *sarama.Config
	debugLog               bool
	db                     db.DBWrapper
	inboxRetry             config.InboxRetry
	groupId                string

	inbox      inbox.InboxMessageService
	subscriber message.Subscriber
	router     *message.Router
}

func NewKafka(config config.Kafka, logger logger.InfraLogger, db db.DBWrapper) Kafka {
	if logger == nil {
		panic("Cannot instantiate a watermill without a logger")
	}

	lg := newWatermillLogger(logger)

	instance := Kafka{
		watermillLogger: lg,
		brokers:         config.Brokers,
		debugLog:        config.DebugLog,
		db:              db,
		groupId:         config.GroupId,
	}

	err := inbox.InboxCheckup(instance.db)
	if err != nil {
		panic("Checking or migration of the InboxMessages failed: " + err.Error())
	}

	instance.inboxRetry.MaxDelay = utiles.ZeroDefault(config.InboxRetry.MaxDelay, 30000)
	instance.inboxRetry.BaseDelay = utiles.ZeroDefault(config.InboxRetry.MaxDelay, 1000)
	instance.inboxRetry.MaxRetries = utiles.ZeroDefault(config.InboxRetry.MaxDelay, 3)
	instance.inboxRetry.ScaleFactor = utiles.ZeroDefault(config.InboxRetry.MaxDelay, 2)

	inboxRepo := repository.NewInboxMessageRepository(instance.db.DB)
	inboxSvc := service.NewInboxMessageService(inboxRepo)
	instance.inbox = inboxSvc

	instance.saramaSubscriberConfig = kafka.DefaultSaramaSubscriberConfig()

	if config.FromBeginning {
		instance.saramaSubscriberConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	err = instance.initiateRouter(config.Router)
	if err != nil {
		panic("Initiating router failed: " + err.Error())
	}

	return instance
}

func (k *Kafka) initiateRouter(routerConfig config.KafkaRouterConfig) error {
	var err error
	var closeTimeOut time.Duration
	if routerConfig.CloseTimeout != "" {
		closeTimeOut, err = time.ParseDuration(routerConfig.CloseTimeout)
		if err != nil {
			k.logger.Warnf("Close time config of the Watermill router config cannot be parsed and set to default, err: %s", err.Error())
		}
	}

	if closeTimeOut == 0 {
		closeTimeOut = defaultRouterCloseTime
	}

	router, err := message.NewRouter(
		message.RouterConfig{
			CloseTimeout: closeTimeOut,
		},
		k.watermillLogger,
	)
	if err != nil {
		panic("Initiating router failed: " + err.Error())
	}

	if !routerConfig.SigTermAndSigIntHandledAlready {
		router.AddPlugin(plugin.SignalsHandler)
	}

	router.AddMiddleware(
		middleware.CorrelationID,
		middleware.Retry{
			MaxRetries:      routerConfig.Retries,
			InitialInterval: time.Millisecond * time.Duration(routerConfig.InitialInterval),
			Logger:          k.watermillLogger,
		}.Middleware,

		middleware.Recoverer,
	)

	subs, err := k.newSubscriber()
	if err != nil {
		return err
	}

	k.subscriber = subs
	k.router = router
	return nil
}
