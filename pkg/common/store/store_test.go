package store

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thealiakbari/todoapp/pkg/common/logger"
)

var (
	redisAddr     = "localhost:6380"
	redisPassword = ""
	redisDB       = 0
)

func newTestStore() Store {
	logger, err := logger.NewInfra("local", "vote-management", "")
	if err != nil {
		panic(err)
	}

	return New(redisAddr, redisPassword, redisDB, logger)
}

func TestPing(t *testing.T) {
	store := newTestStore()

	err := store.Ping(context.Background())
	require.NoError(t, err)
}

func TestSetAndGet(t *testing.T) {
	store := newTestStore()
	key := "testKey"
	type V struct {
		Foo string `json:"foo"`
	}
	value := V{Foo: "bar"}
	ttl := 1 * time.Second

	// Test Set
	err := store.Set(context.Background(), key, value, ttl)
	require.NoError(t, err)

	// Test Get
	var retrievedValue V
	err = store.Get(context.Background(), key, &retrievedValue)
	require.NoError(t, err)
	assert.Equal(t, value, retrievedValue)

	// Optionally wait for TTL to expire
	time.Sleep(ttl + time.Second)

	// Test Get after TTL expiration
	err = store.Get(context.Background(), key, &retrievedValue)
	assert.Error(t, err)
}

func TestDel(t *testing.T) {
	store := newTestStore()
	key := "testKey"
	value := "testValue"

	// Set the key first
	err := store.Set(context.Background(), key, value, 10*time.Second)
	require.NoError(t, err)

	// Test Del
	err = store.Del(context.Background(), key)
	require.NoError(t, err)

	// Test Get after deletion
	var result string
	err = store.Get(context.Background(), key, &result)
	assert.Error(t, err)
}
