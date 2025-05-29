package kafka_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/thealiakbari/hichapp/pkg/common/config"
	"github.com/thealiakbari/hichapp/pkg/common/db"
	"github.com/thealiakbari/hichapp/pkg/common/kafka"
	"github.com/thealiakbari/hichapp/pkg/common/kafka/outbox/domain/enum"
	"github.com/thealiakbari/hichapp/pkg/common/logger"
	"github.com/thealiakbari/hichapp/pkg/common/middleware"
)

var cfg = config.AppConfig{
	DB: config.DB{
		Postgres: config.Postgres{
			Host:          "172.17.0.1",
			Name:          "acms",
			AppName:       "InboxOutboxTest",
			Port:          5432,
			Username:      "postgres",
			Password:      "postgres",
			AutoMigration: false,
		},
	},
	Kafka: config.Kafka{
		Brokers: []string{"localhost:9092"},
		Retry: config.KafkaRetries{
			Retries: 1,
		},
		InboxRetry: config.InboxRetry{},
	},
}

func prettyPrint(note string, in any) {
	bs, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		println(err)
	}
	println(note, string(bs))
}

type something struct {
	Content string `json:"content"`
	Time    string
}

var somethingCreated = kafka.OutboxMessageSchema[something]{
	Name:  "SomethingCreated",
	Topic: "test.something.created.event",
	Type:  enum.OutboxMessageTypeEvent,
}

type inboxOutboxSuite struct {
	suite.Suite
	db     db.DBWrapper
	outbox kafka.OutboxWrapper
}

func (s *inboxOutboxSuite) SetupTest() {
	ctx := context.Background()

	gormDB, err := db.NewPostgresConn(ctx, cfg.DB.Postgres)
	if err != nil {
		panic(err)
	}

	dbw := db.NewDBWrapper(gormDB)
	s.db = dbw
	s.outbox = kafka.NewOutboxWrapper(kafka.OutboxConfig{
		Db: dbw,
	})
}

func (s *inboxOutboxSuite) TestPutAndReceiveMsg() {
	ctxP := context.Background()
	ctx, cancelFn := context.WithCancel(ctxP)

	logger, err := logger.NewInfra("local", "inbox/outbox test", "")
	if err != nil {
		panic(err)
	}
	kafkaRouter := kafka.NewKafka(cfg.Kafka, logger, s.db)
	defer kafkaRouter.Close()

	traceId := uuid.NewString()
	correlationId := uuid.NewString()

	//prettyPrint("Values", map[string]any{
	//	"traceId":       traceId,
	//	"correlationId": correlationId,
	//})

	msg, err := kafka.ToOutboxMessage(
		somethingCreated,
		something{
			Content: "Hello from test codes",
			Time:    time.Now().String(),
		},
	)
	assert.Equal(s.T(), nil, err)

	x, err := s.outbox.Put(ctx, correlationId, msg)
	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), traceId, x.TraceId)
	assert.Equal(s.T(), correlationId, x.AggregateId)

	receivedData := make(chan something, 1)
	kafkaRouter.AddHandler("test_of_something", "test.something.created.event", func(ctx context.Context, p message.Payload) error {
		something := something{}
		err := json.Unmarshal(p, &something)
		if err != nil {
			return err
		}
		receivedData <- something
		return nil
	})

	go kafkaRouter.Run(ctx)

	select {
	case res := <-receivedData:
		assert.Equal(s.T(), "Hello from test codes", res.Content)
		cancelFn()
	case <-time.After(time.Second * 5):
		cancelFn()
		s.T().Errorf("receiving data from kafka timeouts")
	}
}

func (s *inboxOutboxSuite) TestRetry() {
	ctxP := context.Background()
	ctx, cancelFn := context.WithCancel(ctxP)

	logger, err := logger.NewInfra("local", "inbox/outbox test", "")
	if err != nil {
		panic(err)
	}
	kafkaRouter := kafka.NewKafka(cfg.Kafka, logger, s.db)
	defer kafkaRouter.Close()

	traceId := uuid.New()
	correlationId := uuid.NewString()

	ctx = context.WithValue(ctx, middleware.TraceIdKey, traceId)

	//prettyPrint("Values", map[string]any{
	//	"traceId":       traceId,
	//	"correlationId": correlationId,
	//})

	msg, err := kafka.ToOutboxMessage(
		somethingCreated,
		something{
			Content: "Hello from test codes",
			Time:    time.Now().String(),
		},
	)
	assert.Equal(s.T(), nil, err)

	x, err := s.outbox.Put(ctx, correlationId, msg)
	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), traceId.String(), x.TraceId)
	assert.Equal(s.T(), correlationId, x.AggregateId)

	receivedData := make(chan something, 1)

	attempts := 1

	kafkaRouter.AddHandler("test_of_something", "test.something.created.event", func(ctx context.Context, p message.Payload) error {
		if attempts <= 2 {
			attempts += 1
			return kafka.ErrInboxTimeout
		}
		something := something{}
		err := json.Unmarshal(p, &something)
		if err != nil {
			return err
		}
		receivedData <- something
		return nil
	})

	go kafkaRouter.Run(ctx)

	select {
	case res := <-receivedData:
		assert.Equal(s.T(), "Hello from test codes", res.Content)
		assert.Equal(s.T(), attempts, 3)
		cancelFn()
	case <-time.After(time.Second * 5):
		cancelFn()
		s.T().Errorf("receiving data from kafka timeouts")
	}
}

func (s *inboxOutboxSuite) TearDownSuite() {
	db, err := s.db.DB.DB()
	if err != nil {
		panic("failed to sqlDB")
	}

	db.Close()
}

func TestInboxOutbox(t *testing.T) {
	suite.Run(t, new(inboxOutboxSuite))
}
