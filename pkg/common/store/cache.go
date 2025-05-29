package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/thealiakbari/hichapp/pkg/common/logger"
	"go.elastic.co/apm/v2"
)

const (
	AuthAccessTokenKey = "AUTH_ACCESS_TOKEN.SESSION_ID:"
)

type store struct {
	redisClient *redis.Client
	logger      logger.InfraLogger
	_           struct{}
}

type Store interface {
	Ping(ctx context.Context) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration, quiet ...bool) error
	Get(ctx context.Context, key string, model any) error
	Del(ctx context.Context, key string) error
	FindKeys(ctx context.Context, pattern string) (keys []string, err error)
}

func (c *store) Ping(ctx context.Context) error {
	err := c.redisClient.Ping(ctx).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *store) Set(ctx context.Context, key string, value interface{}, ttl time.Duration, quiet ...bool) error {
	// s.injectApmStat(ctx, "Set", key)

	payload, err := json.Marshal(value)
	if err != nil {
		s.logger.Errorf("Can't marshal json: %v", err)
		return err
	}

	err = s.redisClient.Set(ctx, key, payload, ttl).Err()
	if err != nil {
		s.logger.Errorf("Can't set key on redis key: %s : %v", key, err)
		return err
	}

	if len(quiet) > 0 && !quiet[0] {
		s.logger.Debug("Successful set data on redis", logger.String("key", key), logger.Any("value", value))
	}
	return nil
}

func (s *store) Get(ctx context.Context, key string, model any) error {
	// s.injectApmStat(ctx, "Get", key)

	bytes, err := s.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		s.logger.Errorf("Can't get key from redis key: %s : %v", key, err)
		return err
	}

	err = json.Unmarshal(bytes, &model)
	if err != nil {
		s.logger.Errorf("Can't unmarshal json: %v", err)
		return err
	}

	return nil
}

func (s *store) Del(ctx context.Context, key string) error {
	// s.injectApmStat(ctx, "Del", key)

	err := s.redisClient.Del(ctx, key).Err()
	if err != nil {
		s.logger.Errorf("Can't delete key from redis key: %s : %v", key, err)
		return err
	}

	return nil
}

func (s *store) FindKeys(ctx context.Context, pattern string) (keys []string, err error) {
	return s.redisClient.Keys(ctx, pattern).Result()
}

func (s *store) injectApmStat(ctx context.Context, opt string, key string) {
	apmTx := apm.TransactionFromContext(ctx)
	if apmTx == nil {
		s.logger.Debug("No APM transaction found in context for reporting Redis operation")
		return
	}

	stat := apmStat{
		Opt: opt,
		Key: key,
	}
	apmKey := fmt.Sprintf("%s%s", ApmStatKeyPrefix, key)
	apmTx.Context.SetCustom(apmKey, stat)
}

func New(addr string, password string, db int, logger logger.InfraLogger) Store {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	// redisClient.AddHook(newApmHook())

	return &store{
		redisClient: redisClient,
		logger:      logger.ForService(store{}),
	}
}

func GetSessionStoreKey(sessionId string, userId string) string {
	return AuthAccessTokenKey + sessionId + "_" + userId
}
