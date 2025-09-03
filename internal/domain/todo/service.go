package todo

import (
	"context"
	"errors"

	"github.com/thealiakbari/todoapp/internal/domain/todo/entity"
	todoInterface "github.com/thealiakbari/todoapp/internal/ports/inbound/todo"
	"github.com/thealiakbari/todoapp/internal/ports/outbound/todo"
	"github.com/thealiakbari/todoapp/pkg/common/logger"
	appErr "github.com/thealiakbari/todoapp/pkg/common/response"
)

type TodoItemConfig struct {
	Logger       logger.Logger
	TodoItemRepo todo.TodoItemRepository
}

type todoItemService struct {
	TodoItemConfig
}

func NewTodoItemService(config TodoItemConfig) todoInterface.TodoItemService {
	u := todoItemService{config}
	u.Logger = config.Logger.ForService(u)
	return u
}

func (u todoItemService) Create(ctx context.Context, req entity.TodoItem) (res entity.TodoItem, err error) {
	if err = req.Validate(ctx); err != nil {
		u.Logger.Warnf(ctx, "validation error:%v", err)
		return entity.TodoItem{}, &appErr.Error{
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EValidation,
		}
	}

	todoItemEntity, err := u.TodoItemRepo.Create(ctx, req)
	if err != nil {
		u.Logger.Errorf(ctx, "Cannot create todo item: %v", err)
		return entity.TodoItem{}, &appErr.Error{
			ErrCode: 1024,
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EConflict,
		}
	}

	return todoItemEntity, nil
}

func (u todoItemService) Update(ctx context.Context, req entity.TodoItem) (res entity.TodoItem, err error) {
	if err = req.Validate(ctx); err != nil {
		u.Logger.Warnf(ctx, "validation error:%v", err)
		return entity.TodoItem{}, &appErr.Error{
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EValidation,
		}
	}

	err = u.TodoItemRepo.Update(ctx, req)
	if err != nil {
		return entity.TodoItem{}, &appErr.Error{
			ErrCode: 1024,
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EConflict,
		}
	}

	return req, nil
}

func (u todoItemService) GetByIdOrEmpty(ctx context.Context, id string) (res entity.TodoItem, err error) {
	if id == "" {
		err = errors.New("id must not be empty")
		return entity.TodoItem{}, &appErr.Error{
			ErrCode: 1024,
			Cause:   err,
			Message: "Id(TodoItemId) cannot be empty",
			Class:   appErr.EValidation,
		}
	}

	todoItemEntity, err := u.TodoItemRepo.FindByIdOrEmpty(ctx, id)
	if err != nil {
		return entity.TodoItem{}, &appErr.Error{
			ErrCode: 1024,
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EConflict,
		}
	}

	return todoItemEntity, nil
}

func (u todoItemService) Purge(ctx context.Context, id string) (err error) {
	if id == "" {
		err := errors.New("id must not be empty")
		return &appErr.Error{
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EValidation,
		}
	}

	err = u.TodoItemRepo.Purge(ctx, id)
	if err != nil {
		return &appErr.Error{
			ErrCode: 1024,
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EConflict,
		}
	}

	return nil
}

func (u todoItemService) Delete(ctx context.Context, id string) (err error) {
	if id == "" {
		err := errors.New("id must not be empty")
		return &appErr.Error{
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EValidation,
		}
	}

	err = u.TodoItemRepo.Delete(ctx, id)
	if err != nil {
		return &appErr.Error{
			ErrCode: 1024,
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EConflict,
		}
	}

	return nil
}
