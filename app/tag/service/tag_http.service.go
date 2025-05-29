package service

import (
	"github.com/gin-gonic/gin"
	"github.com/thealiakbari/hichapp/app/tag/domain/dto"
	"github.com/thealiakbari/hichapp/app/tag/domain/transform"
	tagService "github.com/thealiakbari/hichapp/internal/tag"
	"github.com/thealiakbari/hichapp/pkg/common/db"
	appErr "github.com/thealiakbari/hichapp/pkg/common/response"
)

type TagHttpApp struct {
	tagSvc tagService.Tag
	db     db.DBWrapper
}

func NewTagHttpApp(tagSvc tagService.Tag, db db.DBWrapper) TagHttpApp {
	return TagHttpApp{
		db:     db,
		tagSvc: tagSvc,
	}
}

// MakeCreate
// @Schemes
// @Summary Create Tag
// @Description This api for create tag
// @Tags tags
// @Accept json
// @Produce json
// @Content-Type application/json
// @Security Bearer
// @Param  body body dto.CreateTagRequest true "Contains information to set data"
// @Success 201  {object}  dto.Tag
// @Failure 400  {object}  appErr.ErrSwaggerResponse
// @Failure 401  {object}  appErr.ErrSwaggerResponse
// @Failure 403  {object}  appErr.ErrSwaggerResponse
// @Failure 422  {object}  appErr.ErrValidationSwaggerResponse
// @Failure 500  {object}  appErr.ErrSwaggerResponse
// @Router /tags [post]
func (t TagHttpApp) MakeCreate() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var req dto.CreateTagRequest
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

		tagEntityResp, err := t.tagSvc.Create(ctx, transform.CreateTagRequestToEntity(req))
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

		appErr.CreatedResponse(ginCtx, transform.TagEntityToTagDto(tagEntityResp))
	}
}

// MakeUpdate
// @Schemes
// @Summary Update Tag
// @Description This api for update tag
// @Tags tags
// @Accept json
// @Produce json
// @Content-Type application/json
// @Security Bearer
// @Param id path string true "Tag Id"
// @Param  body body dto.UpdateTagRequest true "Contains information to set data"
// @Success 200  {object}  dto.Tag
// @Failure 400  {object}  appErr.ErrSwaggerResponse
// @Failure 401  {object}  appErr.ErrSwaggerResponse
// @Failure 403  {object}  appErr.ErrSwaggerResponse
// @Failure 422  {object}  appErr.ErrValidationSwaggerResponse
// @Failure 500  {object}  appErr.ErrSwaggerResponse
// @Router /tags/{id} [put]
func (t TagHttpApp) MakeUpdate() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var req dto.UpdateTagRequest
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

		updateReq, err := transform.UpdateTagRequestToEntity(req, ginCtx.Param("id"))
		if err != nil {
			appErr.HandelError(ginCtx, &appErr.Error{
				Cause:   err,
				Message: err.Error(),
				Class:   appErr.EBadArg,
			})
			return
		}
		tagEntityResp, err := t.tagSvc.Update(ctx, updateReq)
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

		appErr.OKResponse(ginCtx, transform.TagEntityToTagDto(tagEntityResp))
	}
}

// MakeDelete
// @Schemes
// @Summary Delete Tag
// @Description This api for delete tag
// @Tags tags
// @Accept json
// @Produce json
// @Content-Type application/json
// @Security Bearer
// @Param id path string true "Tag Id"
// @Success 204
// @Failure 400  {object}  appErr.ErrSwaggerResponse
// @Failure 401  {object}  appErr.ErrSwaggerResponse
// @Failure 403  {object}  appErr.ErrSwaggerResponse
// @Failure 422  {object}  appErr.ErrValidationSwaggerResponse
// @Failure 500  {object}  appErr.ErrSwaggerResponse
// @Router /tags/{id} [delete]
func (t TagHttpApp) MakeDelete() gin.HandlerFunc {
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

		err = t.tagSvc.Delete(ctx, ginCtx.Param("id"))
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
// @Summary Purge Tag
// @Description This api for purge tag
// @Tags tags
// @Accept json
// @Produce json
// @Content-Type application/json
// @Security Bearer
// @Param id path string true "Tag Id"
// @Success 204
// @Failure 400  {object}  appErr.ErrSwaggerResponse
// @Failure 401  {object}  appErr.ErrSwaggerResponse
// @Failure 403  {object}  appErr.ErrSwaggerResponse
// @Failure 422  {object}  appErr.ErrValidationSwaggerResponse
// @Failure 500  {object}  appErr.ErrSwaggerResponse
// @Router /tags/purge/{id} [delete]
func (t TagHttpApp) MakePurge() gin.HandlerFunc {
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

		err = t.tagSvc.Purge(ctx, ginCtx.Param("id"))
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
// @Summary Get Tag By Id
// @Description This api for tag by id
// @Tags tags
// @Accept json
// @Produce json
// @Content-Type application/json
// @Security Bearer
// @Param id path string true "Tag Id"
// @Success 200  {object} dto.Tag
// @Failure 400  {object}  appErr.ErrSwaggerResponse
// @Failure 401  {object}  appErr.ErrSwaggerResponse
// @Failure 403  {object}  appErr.ErrSwaggerResponse
// @Failure 422  {object}  appErr.ErrValidationSwaggerResponse
// @Failure 500  {object}  appErr.ErrSwaggerResponse
// @Router /tags/{id} [get]
func (t TagHttpApp) MakeGetById() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		tagEntityResp, err := t.tagSvc.GetByIdOrEmpty(ginCtx.Request.Context(), ginCtx.Param("id"))
		if err != nil {
			appErr.HandelError(ginCtx, err)
			return
		}

		appErr.OKResponse(ginCtx, transform.TagEntityToTagDto(tagEntityResp))
	}
}
