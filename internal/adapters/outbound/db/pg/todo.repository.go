package pg

import (
	"context"

	"github.com/thealiakbari/todoapp/internal/domain/todo/entity"
	"github.com/thealiakbari/todoapp/internal/ports/outbound/todo"
	"github.com/thealiakbari/todoapp/pkg/common/db"
)

type todoItemConfig struct {
	db db.DBWrapper
}

func NewTodoItemRepository(db db.DBWrapper) todo.TodoItemRepository {
	return todoItemConfig{
		db: db,
	}
}

func (u todoItemConfig) Create(ctx context.Context, in entity.TodoItem) (res entity.TodoItem, err error) {
	err = db.GormConnection(ctx, u.db.DB).Save(&in).Error
	if err != nil {
		return entity.TodoItem{}, err
	}

	return in, nil
}

func (u todoItemConfig) Update(ctx context.Context, in entity.TodoItem) (err error) {
	err = db.GormConnection(ctx, u.db.DB).Save(&in).Error
	if err != nil {
		return err
	}

	return nil
}

func (u todoItemConfig) FindByIdOrEmpty(ctx context.Context, id string) (res entity.TodoItem, err error) {
	err = db.GormConnection(ctx, u.db.DB).Model(&res).Order("created_at desc").Find(&res, "id = ?", id).Limit(1).Error
	if err != nil {
		return entity.TodoItem{}, err
	}

	return res, nil
}

func (u todoItemConfig) FindByIds(ctx context.Context, ids []string) (res []entity.TodoItem, err error) {
	err = db.GormConnection(ctx, u.db.DB).Model(&res).Find(&res, "id IN (?)", ids).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u todoItemConfig) Purge(ctx context.Context, id string) (err error) {
	err = db.GormConnection(ctx, u.db.DB).Exec("DELETE FROM todo_items WHERE id = ?", id).Error
	if err != nil {
		return err
	}

	return nil
}

func (u todoItemConfig) Delete(ctx context.Context, id string) (err error) {
	err = db.GormConnection(ctx, u.db.DB).Delete(&entity.TodoItem{}, "id = ?", id).Error
	if err != nil {
		return err
	}

	return nil
}

func (u todoItemConfig) FilterFind(ctx context.Context, query []any, order string, limit int, offset int) (res []entity.TodoItem, err error) {
	err = db.GormConnection(ctx, u.db.DB).Model(&res).Order(order).
		Limit(limit).
		Offset(offset).
		Find(&res, query...).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u todoItemConfig) FilterCount(ctx context.Context, query []any) (res int64, err error) {
	countQuery := db.GormConnection(ctx, u.db.DB).Model(&entity.TodoItem{})
	if len(query) > 1 {
		countQuery = countQuery.Where(query[0], query[1:]...)
	} else if len(query) == 1 {
		countQuery = countQuery.Where(query[0])
	}

	err = countQuery.Count(&res).Error
	if err != nil {
		return 0, err
	}

	return res, nil
}
