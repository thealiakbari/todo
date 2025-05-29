package service

import (
	"github.com/gin-gonic/gin"
	"github.com/thealiakbari/hichapp/app/user/domain/dto"
	"github.com/thealiakbari/hichapp/app/user/domain/transform"
	userService "github.com/thealiakbari/hichapp/internal/user"
	"github.com/thealiakbari/hichapp/pkg/common/db"
	appErr "github.com/thealiakbari/hichapp/pkg/common/response"
)

type UserHttpApp struct {
	userSvc userService.User
	db      db.DBWrapper
}

func NewUserHttpApp(userSvc userService.User, db db.DBWrapper) UserHttpApp {
	return UserHttpApp{
		db:      db,
		userSvc: userSvc,
	}
}

// MakeCreate
// @Schemes
// @Summary Create User
// @Description This api for create vote
// @Tags users
// @Accept json
// @Produce json
// @Content-Type application/json
// @Security Bearer
// @Param  body body dto.CreateUserRequest true "Contains information to set data"
// @Success 201  {object}  dto.User
// @Failure 400  {object}  appErr.ErrSwaggerResponse
// @Failure 401  {object}  appErr.ErrSwaggerResponse
// @Failure 403  {object}  appErr.ErrSwaggerResponse
// @Failure 422  {object}  appErr.ErrValidationSwaggerResponse
// @Failure 500  {object}  appErr.ErrSwaggerResponse
// @Router /users [post]
func (t UserHttpApp) MakeCreate() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var req dto.CreateUserRequest
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

		userEntityResp, err := t.userSvc.Create(ctx, transform.CreateUserRequestToEntity(req))
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

		appErr.CreatedResponse(ginCtx, transform.UserEntityToUserDto(userEntityResp))
	}
}

// MakeUpdate
// @Schemes
// @Summary Update User
// @Description This api for update vote
// @Tags users
// @Accept json
// @Produce json
// @Content-Type application/json
// @Security Bearer
// @Param id path string true "User Id"
// @Param  body body dto.UpdateUserRequest true "Contains information to set data"
// @Success 200  {object}  dto.User
// @Failure 400  {object}  appErr.ErrSwaggerResponse
// @Failure 401  {object}  appErr.ErrSwaggerResponse
// @Failure 403  {object}  appErr.ErrSwaggerResponse
// @Failure 422  {object}  appErr.ErrValidationSwaggerResponse
// @Failure 500  {object}  appErr.ErrSwaggerResponse
// @Router /users/{id} [put]
func (t UserHttpApp) MakeUpdate() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var req dto.UpdateUserRequest
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

		updateReq, err := transform.UpdateUserRequestToEntity(req, ginCtx.Param("id"))
		if err != nil {
			appErr.HandelError(ginCtx, &appErr.Error{
				Cause:   err,
				Message: err.Error(),
				Class:   appErr.EBadArg,
			})
			return
		}
		userEntityResp, err := t.userSvc.Update(ctx, updateReq)
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

		appErr.OKResponse(ginCtx, transform.UserEntityToUserDto(userEntityResp))
	}
}

// MakeDelete
// @Schemes
// @Summary Delete User
// @Description This api for delete vote
// @Tags users
// @Accept json
// @Produce json
// @Content-Type application/json
// @Security Bearer
// @Param id path string true "User Id"
// @Success 204
// @Failure 400  {object}  appErr.ErrSwaggerResponse
// @Failure 401  {object}  appErr.ErrSwaggerResponse
// @Failure 403  {object}  appErr.ErrSwaggerResponse
// @Failure 422  {object}  appErr.ErrValidationSwaggerResponse
// @Failure 500  {object}  appErr.ErrSwaggerResponse
// @Router /users/{id} [delete]
func (t UserHttpApp) MakeDelete() gin.HandlerFunc {
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

		err = t.userSvc.Delete(ctx, ginCtx.Param("id"))
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
// @Summary Purge User
// @Description This api for purge vote
// @Tags users
// @Accept json
// @Produce json
// @Content-Type application/json
// @Security Bearer
// @Param id path string true "User Id"
// @Success 204
// @Failure 400  {object}  appErr.ErrSwaggerResponse
// @Failure 401  {object}  appErr.ErrSwaggerResponse
// @Failure 403  {object}  appErr.ErrSwaggerResponse
// @Failure 422  {object}  appErr.ErrValidationSwaggerResponse
// @Failure 500  {object}  appErr.ErrSwaggerResponse
// @Router /users/purge/{id} [delete]
func (t UserHttpApp) MakePurge() gin.HandlerFunc {
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

		err = t.userSvc.Purge(ctx, ginCtx.Param("id"))
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
// @Summary Get User By Id
// @Description This api for vote by id
// @Tags users
// @Accept json
// @Produce json
// @Content-Type application/json
// @Security Bearer
// @Param id path string true "User Id"
// @Success 200  {object} dto.User
// @Failure 400  {object}  appErr.ErrSwaggerResponse
// @Failure 401  {object}  appErr.ErrSwaggerResponse
// @Failure 403  {object}  appErr.ErrSwaggerResponse
// @Failure 422  {object}  appErr.ErrValidationSwaggerResponse
// @Failure 500  {object}  appErr.ErrSwaggerResponse
// @Router /users/{id} [get]
func (t UserHttpApp) MakeGetById() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		userEntityResp, err := t.userSvc.GetByIdOrEmpty(ginCtx.Request.Context(), ginCtx.Param("id"))
		if err != nil {
			appErr.HandelError(ginCtx, err)
			return
		}

		appErr.OKResponse(ginCtx, transform.UserEntityToUserDto(userEntityResp))
	}
}
