package service

import (
	"github.com/gin-gonic/gin"
	"github.com/thealiakbari/hichapp/app/poll/domain/dto"
	"github.com/thealiakbari/hichapp/app/poll/domain/transform"
	pollService "github.com/thealiakbari/hichapp/internal/poll"
	"github.com/thealiakbari/hichapp/pkg/common/db"
	appErr "github.com/thealiakbari/hichapp/pkg/common/response"
)

type PollHttpApp struct {
	pollSvc pollService.Poll
	db      db.DBWrapper
}

func NewPollHttpApp(pollSvc pollService.Poll, db db.DBWrapper) PollHttpApp {
	return PollHttpApp{
		db:      db,
		pollSvc: pollSvc,
	}
}

// MakeCreate
// @Schemes
// @Summary Create Poll
// @Description This api for create poll
// @Tags polls
// @Accept json
// @Produce json
// @Content-Type application/json
// @Security Bearer
// @Param  body body dto.CreatePollRequest true "Contains information to set data"
// @Success 201  {object}  dto.Poll
// @Failure 400  {object}  appErr.ErrSwaggerResponse
// @Failure 401  {object}  appErr.ErrSwaggerResponse
// @Failure 403  {object}  appErr.ErrSwaggerResponse
// @Failure 422  {object}  appErr.ErrValidationSwaggerResponse
// @Failure 500  {object}  appErr.ErrSwaggerResponse
// @Router /polls [post]
func (t PollHttpApp) MakeCreate() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var req dto.CreatePollRequest
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

		pollEntityResp, err := t.pollSvc.Create(ctx, transform.CreatePollRequestToEntity(req))
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

		appErr.CreatedResponse(ginCtx, transform.PollEntityToPollDto(pollEntityResp))
	}
}

// MakeUpdate
// @Schemes
// @Summary Update Poll
// @Description This api for update poll
// @Tags polls
// @Accept json
// @Produce json
// @Content-Type application/json
// @Security Bearer
// @Param id path string true "Poll Id"
// @Param  body body dto.UpdatePollRequest true "Contains information to set data"
// @Success 200  {object}  dto.Poll
// @Failure 400  {object}  appErr.ErrSwaggerResponse
// @Failure 401  {object}  appErr.ErrSwaggerResponse
// @Failure 403  {object}  appErr.ErrSwaggerResponse
// @Failure 422  {object}  appErr.ErrValidationSwaggerResponse
// @Failure 500  {object}  appErr.ErrSwaggerResponse
// @Router /polls/{id} [put]
func (t PollHttpApp) MakeUpdate() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var req dto.UpdatePollRequest
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

		updateReq, err := transform.UpdatePollRequestToEntity(req, ginCtx.Param("id"))
		if err != nil {
			appErr.HandelError(ginCtx, &appErr.Error{
				Cause:   err,
				Message: err.Error(),
				Class:   appErr.EBadArg,
			})
			return
		}
		pollEntityResp, err := t.pollSvc.Update(ctx, updateReq)
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

		appErr.OKResponse(ginCtx, transform.PollEntityToPollDto(pollEntityResp))
	}
}

// MakeDelete
// @Schemes
// @Summary Delete Poll
// @Description This api for delete poll
// @Tags polls
// @Accept json
// @Produce json
// @Content-Type application/json
// @Security Bearer
// @Param id path string true "Poll Id"
// @Success 204
// @Failure 400  {object}  appErr.ErrSwaggerResponse
// @Failure 401  {object}  appErr.ErrSwaggerResponse
// @Failure 403  {object}  appErr.ErrSwaggerResponse
// @Failure 422  {object}  appErr.ErrValidationSwaggerResponse
// @Failure 500  {object}  appErr.ErrSwaggerResponse
// @Router /polls/{id} [delete]
func (t PollHttpApp) MakeDelete() gin.HandlerFunc {
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

		err = t.pollSvc.Delete(ctx, ginCtx.Param("id"))
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
// @Summary Purge Poll
// @Description This api for purge poll
// @Tags polls
// @Accept json
// @Produce json
// @Content-Type application/json
// @Security Bearer
// @Param id path string true "Poll Id"
// @Success 204
// @Failure 400  {object}  appErr.ErrSwaggerResponse
// @Failure 401  {object}  appErr.ErrSwaggerResponse
// @Failure 403  {object}  appErr.ErrSwaggerResponse
// @Failure 422  {object}  appErr.ErrValidationSwaggerResponse
// @Failure 500  {object}  appErr.ErrSwaggerResponse
// @Router /polls/purge/{id} [delete]
func (t PollHttpApp) MakePurge() gin.HandlerFunc {
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

		err = t.pollSvc.Purge(ctx, ginCtx.Param("id"))
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
// @Summary Get Poll By Id
// @Description This api for poll by id
// @Tags polls
// @Accept json
// @Produce json
// @Content-Type application/json
// @Security Bearer
// @Param id path string true "Poll Id"
// @Success 200  {object} dto.Poll
// @Failure 400  {object}  appErr.ErrSwaggerResponse
// @Failure 401  {object}  appErr.ErrSwaggerResponse
// @Failure 403  {object}  appErr.ErrSwaggerResponse
// @Failure 422  {object}  appErr.ErrValidationSwaggerResponse
// @Failure 500  {object}  appErr.ErrSwaggerResponse
// @Router /polls/{id} [get]
func (t PollHttpApp) MakeGetById() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		pollEntityResp, err := t.pollSvc.GetByIdOrEmpty(ginCtx.Request.Context(), ginCtx.Param("id"))
		if err != nil {
			appErr.HandelError(ginCtx, err)
			return
		}

		appErr.OKResponse(ginCtx, transform.PollEntityToPollDto(pollEntityResp))
	}
}
