package todo

import (
	"context"
	"github.com/thealiakbari/todoapp/internal/domain/todo/entity"
)

type TodoItemService interface {
	Create(ctx context.Context, entity entity.TodoItem) (res entity.TodoItem, err error)
	Update(ctx context.Context, entity entity.TodoItem) (res entity.TodoItem, err error)
	GetByIdOrEmpty(ctx context.Context, id string) (res entity.TodoItem, err error)
	Delete(ctx context.Context, id string) (err error)
	Purge(ctx context.Context, id string) (err error)
}
