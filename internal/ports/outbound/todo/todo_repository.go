package todo

import (
	"context"
	"github.com/thealiakbari/todoapp/internal/domain/todo/entity"
)

type TodoItemRepository interface {
	Create(ctx context.Context, in entity.TodoItem) (res entity.TodoItem, err error)
	Update(ctx context.Context, in entity.TodoItem) (err error)
	FindByIds(ctx context.Context, ids []string) (res []entity.TodoItem, err error)
	FindByIdOrEmpty(ctx context.Context, id string) (res entity.TodoItem, err error)
	Purge(ctx context.Context, id string) (err error)
	Delete(ctx context.Context, id string) (err error)
	FilterFind(ctx context.Context, query []any, order string, limit int, offset int) (res []entity.TodoItem, err error)
	FilterCount(ctx context.Context, query []any) (res int64, err error)
}
