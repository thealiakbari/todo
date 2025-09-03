package service

import (
	"github.com/gin-gonic/gin"
	"github.com/thealiakbari/todoapp/internal/application/todo/domain/dto"
	"github.com/thealiakbari/todoapp/internal/application/todo/domain/transform"
	todoInterface "github.com/thealiakbari/todoapp/internal/ports/inbound/todo"
	"github.com/thealiakbari/todoapp/pkg/common/db"
	appErr "github.com/thealiakbari/todoapp/pkg/common/response"
)

type TodoItemHttpApp struct {
	todoItemSvc todoInterface.TodoItemService
	db          db.DBWrapper
}

func NewTodoItemHttpApp(todoItemSvc todoInterface.TodoItemService, db db.DBWrapper) TodoItemHttpApp {
	return TodoItemHttpApp{
		db:          db,
		todoItemSvc: todoItemSvc,
	}
}

// MakeCreate
// @Schemes
// @Summary Create TodoItem
// @Description This api for create poll
// @Tags todo-items
// @Accept json
// @Produce json
// @Content-Type application/json
// @Security Bearer
// @Param  body body dto.CreateTodoItemRequest true "Contains information to set data"
// @Success 201  {object}  dto.TodoItem
// @Failure 400  {object}  appErr.ErrSwaggerResponse
// @Failure 401  {object}  appErr.ErrSwaggerResponse
// @Failure 403  {object}  appErr.ErrSwaggerResponse
// @Failure 422  {object}  appErr.ErrValidationSwaggerResponse
// @Failure 500  {object}  appErr.ErrSwaggerResponse
// @Router /todo-items [post]
func (t TodoItemHttpApp) MakeCreate() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var req dto.CreateTodoItemRequest
		if err := ginCtx.ShouldBindJSON(&req); err != nil {
			appErr.HandelError(ginCtx, &appErr.Error{
				Cause:   err,
				Message: err.Error(),
				Class:   appErr.EBadArg,
			})
			return
		}

		if err := req.Validate(ginCtx.Request.Context()); err != nil {
			appErr.HandelError(ginCtx, &appErr.Error{
				Cause:   err,
				Message: err.Error(),
				Class:   appErr.EValidation,
			})
			return
		}

		tx, ctx, err := db.BeginTx(ginCtx.Request.Context(), t.db.DB)
		if err != nil {
			appErr.HandelError(ginCtx, &appErr.Error{
				Cause:   err,
				Message: err.Error(),
				Class:   appErr.EConflict,
			})
			return
		}

		defer func() {
			if err != nil {
				if err := tx.Rollback().Error; err != nil {
					appErr.HandelError(ginCtx, &appErr.Error{
						Cause:   err,
						Message: err.Error(),
						Class:   appErr.EConflict,
					})
					return
				}
				appErr.HandelError(ginCtx, err)
				return
			}
		}()

		pollEntityResp, err := t.todoItemSvc.Create(ctx, transform.CreateTodoItemRequestToEntity(req))
		if err != nil {
			appErr.HandelError(ginCtx, err)
			return
		}

		if err = tx.Commit().Error; err != nil {
			appErr.HandelError(ginCtx, &appErr.Error{
				Cause:   err,
				Message: err.Error(),
				Class:   appErr.EConflict,
			})
			return
		}

		appErr.CreatedResponse(ginCtx, transform.TodoItemEntityToTodoItemDto(pollEntityResp))
	}
}

// MakeUpdate
// @Schemes
// @Summary Update TodoItem
// @Description This api for update poll
// @Tags todo-items
// @Accept json
// @Produce json
// @Content-Type application/json
// @Security Bearer
// @Param id path string true "TodoItem Id"
// @Param  body body dto.UpdateTodoItemRequest true "Contains information to set data"
// @Success 200  {object}  dto.TodoItem
// @Failure 400  {object}  appErr.ErrSwaggerResponse
// @Failure 401  {object}  appErr.ErrSwaggerResponse
// @Failure 403  {object}  appErr.ErrSwaggerResponse
// @Failure 422  {object}  appErr.ErrValidationSwaggerResponse
// @Failure 500  {object}  appErr.ErrSwaggerResponse
// @Router /todo-items/{id} [put]
func (t TodoItemHttpApp) MakeUpdate() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var req dto.UpdateTodoItemRequest
		if err := ginCtx.ShouldBindJSON(&req); err != nil {
			appErr.HandelError(ginCtx, &appErr.Error{
				Cause:   err,
				Message: err.Error(),
				Class:   appErr.EBadArg,
			})
			return
		}

		tx, ctx, err := db.BeginTx(ginCtx.Request.Context(), t.db.DB)
		if err != nil {
			appErr.HandelError(ginCtx, &appErr.Error{
				Cause:   err,
				Message: err.Error(),
				Class:   appErr.EConflict,
			})
			return
		}

		defer func() {
			if err != nil {
				if err := tx.Rollback().Error; err != nil {
					appErr.HandelError(ginCtx, &appErr.Error{
						Cause:   err,
						Message: err.Error(),
						Class:   appErr.EConflict,
					})
					return
				}
				appErr.HandelError(ginCtx, err)
				return
			}
		}()

		updateReq, err := transform.UpdateTodoItemRequestToEntity(req, ginCtx.Param("id"))
		if err != nil {
			appErr.HandelError(ginCtx, &appErr.Error{
				Cause:   err,
				Message: err.Error(),
				Class:   appErr.EBadArg,
			})
			return
		}
		pollEntityResp, err := t.todoItemSvc.Update(ctx, updateReq)
		if err != nil {
			appErr.HandelError(ginCtx, err)
			return
		}

		if err = tx.Commit().Error; err != nil {
			appErr.HandelError(ginCtx, &appErr.Error{
				Cause:   err,
				Message: err.Error(),
				Class:   appErr.EConflict,
			})
			return
		}

		appErr.OKResponse(ginCtx, transform.TodoItemEntityToTodoItemDto(pollEntityResp))
	}
}

// MakeDelete
// @Schemes
// @Summary Delete TodoItem
// @Description This api for delete poll
// @Tags todo-items
// @Accept json
// @Produce json
// @Content-Type application/json
// @Security Bearer
// @Param id path string true "TodoItem Id"
// @Success 204
// @Failure 400  {object}  appErr.ErrSwaggerResponse
// @Failure 401  {object}  appErr.ErrSwaggerResponse
// @Failure 403  {object}  appErr.ErrSwaggerResponse
// @Failure 422  {object}  appErr.ErrValidationSwaggerResponse
// @Failure 500  {object}  appErr.ErrSwaggerResponse
// @Router /todo-items/{id} [delete]
func (t TodoItemHttpApp) MakeDelete() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		tx, ctx, err := db.BeginTx(ginCtx.Request.Context(), t.db.DB)
		if err != nil {
			appErr.HandelError(ginCtx, &appErr.Error{
				Cause:   err,
				Message: err.Error(),
				Class:   appErr.EConflict,
			})
			return
		}

		defer func() {
			if err != nil {
				if err := tx.Rollback().Error; err != nil {
					appErr.HandelError(ginCtx, &appErr.Error{
						Cause:   err,
						Message: err.Error(),
						Class:   appErr.EConflict,
					})
					return
				}
				appErr.HandelError(ginCtx, err)
				return
			}
		}()

		err = t.todoItemSvc.Delete(ctx, ginCtx.Param("id"))
		if err != nil {
			appErr.HandelError(ginCtx, err)
			return
		}

		if err = tx.Commit().Error; err != nil {
			appErr.HandelError(ginCtx, &appErr.Error{
				Cause:   err,
				Message: err.Error(),
				Class:   appErr.EConflict,
			})
			return
		}

		appErr.NoContentResponse(ginCtx)
	}
}

// MakePurge
// @Schemes
// @Summary Purge TodoItem
// @Description This api for purge poll
// @Tags todo-items
// @Accept json
// @Produce json
// @Content-Type application/json
// @Security Bearer
// @Param id path string true "TodoItem Id"
// @Success 204
// @Failure 400  {object}  appErr.ErrSwaggerResponse
// @Failure 401  {object}  appErr.ErrSwaggerResponse
// @Failure 403  {object}  appErr.ErrSwaggerResponse
// @Failure 422  {object}  appErr.ErrValidationSwaggerResponse
// @Failure 500  {object}  appErr.ErrSwaggerResponse
// @Router /todo-items/purge/{id} [delete]
func (t TodoItemHttpApp) MakePurge() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		tx, ctx, err := db.BeginTx(ginCtx.Request.Context(), t.db.DB)
		if err != nil {
			appErr.HandelError(ginCtx, &appErr.Error{
				Cause:   err,
				Message: err.Error(),
				Class:   appErr.EConflict,
			})
			return
		}

		defer func() {
			if err != nil {
				if err := tx.Rollback().Error; err != nil {
					appErr.HandelError(ginCtx, &appErr.Error{
						Cause:   err,
						Message: err.Error(),
						Class:   appErr.EConflict,
					})
					return
				}
			}
		}()

		err = t.todoItemSvc.Purge(ctx, ginCtx.Param("id"))
		if err != nil {
			appErr.HandelError(ginCtx, err)
			return
		}

		if err = tx.Commit().Error; err != nil {
			appErr.HandelError(ginCtx, &appErr.Error{
				Cause:   err,
				Message: err.Error(),
				Class:   appErr.EConflict,
			})
			return
		}

		appErr.NoContentResponse(ginCtx)
	}
}

// MakeGetById
// @Schemes
// @Summary Get TodoItem By Id
// @Description This api for poll by id
// @Tags todo-items
// @Accept json
// @Produce json
// @Content-Type application/json
// @Security Bearer
// @Param id path string true "TodoItem Id"
// @Success 200  {object} dto.TodoItem
// @Failure 400  {object}  appErr.ErrSwaggerResponse
// @Failure 401  {object}  appErr.ErrSwaggerResponse
// @Failure 403  {object}  appErr.ErrSwaggerResponse
// @Failure 422  {object}  appErr.ErrValidationSwaggerResponse
// @Failure 500  {object}  appErr.ErrSwaggerResponse
// @Router /todo-items/{id} [get]
func (t TodoItemHttpApp) MakeGetById() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		pollEntityResp, err := t.todoItemSvc.GetByIdOrEmpty(ginCtx.Request.Context(), ginCtx.Param("id"))
		if err != nil {
			appErr.HandelError(ginCtx, err)
			return
		}

		appErr.OKResponse(ginCtx, transform.TodoItemEntityToTodoItemDto(pollEntityResp))
	}
}
