package tag

import (
	"context"
	"errors"

	"github.com/thealiakbari/hichapp/internal/tag/domain/entity"
	tagRepo "github.com/thealiakbari/hichapp/internal/tag/domain/repository"
	"github.com/thealiakbari/hichapp/pkg/common/logger"
	appErr "github.com/thealiakbari/hichapp/pkg/common/response"
)

type Tag interface {
	Create(ctx context.Context, entity entity.Tag) (res entity.Tag, err error)
	Update(ctx context.Context, entity entity.Tag) (res entity.Tag, err error)
	GetByIdOrEmpty(ctx context.Context, id string) (res entity.Tag, err error)
	Delete(ctx context.Context, id string) (err error)
	Purge(ctx context.Context, id string) (err error)
}

type TagConfig struct {
	Logger  logger.Logger
	TagRepo tagRepo.TagRepository
}

type tagService struct {
	TagConfig
}

func NewTagService(config TagConfig) Tag {
	u := tagService{config}
	u.Logger = config.Logger.ForService(u)
	return u
}

func (u tagService) Create(ctx context.Context, req entity.Tag) (res entity.Tag, err error) {
	if err = req.Validate(ctx); err != nil {
		return entity.Tag{}, &appErr.Error{
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EValidation,
		}
	}

	tagEntity, err := u.TagRepo.Create(ctx, req)
	if err != nil {
		return entity.Tag{}, &appErr.Error{
			ErrCode: 1024,
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EConflict,
		}
	}

	return tagEntity, nil
}

func (u tagService) Update(ctx context.Context, req entity.Tag) (res entity.Tag, err error) {
	if err = req.Validate(ctx); err != nil {
		return entity.Tag{}, &appErr.Error{
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EValidation,
		}
	}

	err = u.TagRepo.Update(ctx, req)
	if err != nil {
		return entity.Tag{}, &appErr.Error{
			ErrCode: 1024,
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EConflict,
		}
	}

	return req, nil
}

func (u tagService) GetByIdOrEmpty(ctx context.Context, id string) (res entity.Tag, err error) {
	if id == "" {
		err := errors.New("id must not be empty")
		return entity.Tag{}, &appErr.Error{
			ErrCode: 1024, // todo
			Cause:   err,
			Message: "Id(TagId) cannot be empty",
			Class:   appErr.EValidation,
		}
	}

	tagEntity, err := u.TagRepo.FindByIdOrEmpty(ctx, id)
	if err != nil {
		return entity.Tag{}, &appErr.Error{
			ErrCode: 1024, // todo
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EConflict,
		}
	}

	return tagEntity, nil
}

func (u tagService) Purge(ctx context.Context, id string) (err error) {
	if id == "" {
		err := errors.New("id must not be empty")
		return &appErr.Error{
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EValidation,
		}
	}

	err = u.TagRepo.Purge(ctx, id)
	if err != nil {
		return &appErr.Error{
			ErrCode: 1024, // todo
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EConflict,
		}
	}

	return nil
}

func (u tagService) Delete(ctx context.Context, id string) (err error) {
	if id == "" {
		err := errors.New("id must not be empty")
		return &appErr.Error{
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EValidation,
		}
	}

	err = u.TagRepo.Delete(ctx, id)
	if err != nil {
		return &appErr.Error{
			ErrCode: 1024, // todo
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EConflict,
		}
	}

	return nil
}
