package poll

import (
	"context"
	"errors"

	"github.com/thealiakbari/hichapp/internal/poll/domain/entity"
	pollRepo "github.com/thealiakbari/hichapp/internal/poll/domain/repository"
	"github.com/thealiakbari/hichapp/pkg/common/logger"
	appErr "github.com/thealiakbari/hichapp/pkg/common/response"
)

type Poll interface {
	Create(ctx context.Context, entity entity.Poll) (res entity.Poll, err error)
	Update(ctx context.Context, entity entity.Poll) (res entity.Poll, err error)
	GetByIdOrEmpty(ctx context.Context, id string) (res entity.Poll, err error)
	Delete(ctx context.Context, id string) (err error)
	Purge(ctx context.Context, id string) (err error)
}

type PollConfig struct {
	Logger   logger.Logger
	PollRepo pollRepo.PollRepository
}

type pollService struct {
	PollConfig
}

func NewPollService(config PollConfig) Poll {
	u := pollService{config}
	u.Logger = config.Logger.ForService(u)
	return u
}

func (u pollService) Create(ctx context.Context, req entity.Poll) (res entity.Poll, err error) {
	if err = req.Validate(ctx); err != nil {
		return entity.Poll{}, &appErr.Error{
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EValidation,
		}
	}

	pollEntity, err := u.PollRepo.Create(ctx, req)
	if err != nil {
		return entity.Poll{}, &appErr.Error{
			ErrCode: 1024,
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EConflict,
		}
	}

	return pollEntity, nil
}

func (u pollService) Update(ctx context.Context, req entity.Poll) (res entity.Poll, err error) {
	if err = req.Validate(ctx); err != nil {
		return entity.Poll{}, &appErr.Error{
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EValidation,
		}
	}

	err = u.PollRepo.Update(ctx, req)
	if err != nil {
		return entity.Poll{}, &appErr.Error{
			ErrCode: 1024,
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EConflict,
		}
	}

	return req, nil
}

func (u pollService) GetByIdOrEmpty(ctx context.Context, id string) (res entity.Poll, err error) {
	if id == "" {
		err := errors.New("id must not be empty")
		return entity.Poll{}, &appErr.Error{
			ErrCode: 1024, // todo
			Cause:   err,
			Message: "Id(PollId) cannot be empty",
			Class:   appErr.EValidation,
		}
	}

	pollEntity, err := u.PollRepo.FindByIdOrEmpty(ctx, id)
	if err != nil {
		return entity.Poll{}, &appErr.Error{
			ErrCode: 1024, // todo
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EConflict,
		}
	}

	return pollEntity, nil
}

func (u pollService) Purge(ctx context.Context, id string) (err error) {
	if id == "" {
		err := errors.New("id must not be empty")
		return &appErr.Error{
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EValidation,
		}
	}

	err = u.PollRepo.Purge(ctx, id)
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

func (u pollService) Delete(ctx context.Context, id string) (err error) {
	if id == "" {
		err := errors.New("id must not be empty")
		return &appErr.Error{
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EValidation,
		}
	}

	err = u.PollRepo.Delete(ctx, id)
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
