package todo

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thealiakbari/todoapp/internal/domain/todo/entity"
	"github.com/thealiakbari/todoapp/pkg/common/logger"
	appErr "github.com/thealiakbari/todoapp/pkg/common/response"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) Create(ctx context.Context, in entity.TodoItem) (entity.TodoItem, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(entity.TodoItem), args.Error(1)
}

func (m *mockRepo) Update(ctx context.Context, in entity.TodoItem) error {
	args := m.Called(ctx, in)
	return args.Error(0)
}

func (m *mockRepo) FindByIds(ctx context.Context, ids []string) ([]entity.TodoItem, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]entity.TodoItem), args.Error(1)
}

func (m *mockRepo) FindByIdOrEmpty(ctx context.Context, id string) (entity.TodoItem, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.TodoItem), args.Error(1)
}

func (m *mockRepo) Purge(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockRepo) FilterFind(ctx context.Context, query []any, order string, limit int, offset int) ([]entity.TodoItem, error) {
	args := m.Called(ctx, query, order, limit, offset)
	return args.Get(0).([]entity.TodoItem), args.Error(1)
}

func (m *mockRepo) FilterCount(ctx context.Context, query []any) (int64, error) {
	args := m.Called(ctx, query)
	return args.Get(0).(int64), args.Error(1)
}

func TestCreate_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	log, err := logger.New(
		"local",
		"todoapp",
		"todoapp",
	)
	service := NewTodoItemService(TodoItemConfig{
		Logger:       log,
		TodoItemRepo: repo,
	})

	item := entity.TodoItem{Description: "test", DueDate: "2025-01-01"}
	repo.On("Create", ctx, item).Return(item, nil)

	res, err := service.Create(ctx, item)
	assert.NoError(t, err)
	assert.Equal(t, item, res)
	repo.AssertExpectations(t)
}

func TestCreate_ValidationError(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	log, err := logger.New(
		"local",
		"todoapp",
		"todoapp",
	)
	service := NewTodoItemService(TodoItemConfig{
		Logger:       log,
		TodoItemRepo: repo,
	})

	item := entity.TodoItem{}

	res, err := service.Create(ctx, item)
	assert.Error(t, err)
	assert.IsType(t, &appErr.Error{}, err)
	assert.Equal(t, entity.TodoItem{}, res)
}

func TestCreate_RepoError(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	log, err := logger.New(
		"local",
		"todoapp",
		"todoapp",
	)
	service := NewTodoItemService(TodoItemConfig{
		Logger:       log,
		TodoItemRepo: repo,
	})

	item := entity.TodoItem{Description: "test", DueDate: "2025-01-01"}
	repo.On("Create", ctx, item).Return(entity.TodoItem{}, errors.New("db error"))

	res, err := service.Create(ctx, item)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
	assert.Equal(t, entity.TodoItem{}, res)
}

func TestGetByIdOrEmpty_EmptyId(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	log, err := logger.New(
		"local",
		"todoapp",
		"todoapp",
	)
	service := NewTodoItemService(TodoItemConfig{
		Logger:       log,
		TodoItemRepo: repo,
	})

	res, err := service.GetByIdOrEmpty(ctx, "")
	assert.Error(t, err)
	assert.Equal(t, entity.TodoItem{}, res)
	assert.IsType(t, &appErr.Error{}, err)
}

func TestGetByIdOrEmpty_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	log, err := logger.New(
		"local",
		"todoapp",
		"todoapp",
	)
	service := NewTodoItemService(TodoItemConfig{
		Logger:       log,
		TodoItemRepo: repo,
	})

	expected := entity.TodoItem{Description: "test", DueDate: "2025-01-01"}
	repo.On("FindByIdOrEmpty", ctx, "123").Return(expected, nil)

	res, err := service.GetByIdOrEmpty(ctx, "123")
	assert.NoError(t, err)
	assert.Equal(t, expected, res)
}

func TestDelete_EmptyId(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	log, err := logger.New(
		"local",
		"todoapp",
		"todoapp",
	)
	service := NewTodoItemService(TodoItemConfig{
		Logger:       log,
		TodoItemRepo: repo,
	})

	err = service.Delete(ctx, "")
	assert.Error(t, err)
	assert.IsType(t, &appErr.Error{}, err)
}

func TestDelete_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	log, err := logger.New(
		"local",
		"todoapp",
		"todoapp",
	)
	service := NewTodoItemService(TodoItemConfig{
		Logger:       log,
		TodoItemRepo: repo,
	})

	repo.On("Delete", ctx, "123").Return(nil)

	err = service.Delete(ctx, "123")
	assert.NoError(t, err)
}
