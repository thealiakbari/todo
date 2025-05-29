package user

import (
	"context"
	"errors"

	"github.com/thealiakbari/hichapp/internal/user/domain/entity"
	userRepo "github.com/thealiakbari/hichapp/internal/user/domain/repository"
	"github.com/thealiakbari/hichapp/pkg/common/logger"
	appErr "github.com/thealiakbari/hichapp/pkg/common/response"
)

type User interface {
	Create(ctx context.Context, entity entity.User) (res entity.User, err error)
	Update(ctx context.Context, entity entity.User) (res entity.User, err error)
	GetByIdOrEmpty(ctx context.Context, id string) (res entity.User, err error)
	Delete(ctx context.Context, id string) (err error)
	Purge(ctx context.Context, id string) (err error)
}

type UserConfig struct {
	Logger   logger.Logger
	UserRepo userRepo.UserRepository
}

type userService struct {
	UserConfig
}

func NewUserService(config UserConfig) User {
	u := userService{config}
	u.Logger = config.Logger.ForService(u)
	return u
}

func (u userService) Create(ctx context.Context, req entity.User) (res entity.User, err error) {
	if err = req.Validate(ctx); err != nil {
		return entity.User{}, &appErr.Error{
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EValidation,
		}
	}

	userEntity, err := u.UserRepo.Create(ctx, req)
	if err != nil {
		return entity.User{}, &appErr.Error{
			ErrCode: 1024,
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EConflict,
		}
	}

	return userEntity, nil
}

func (u userService) Update(ctx context.Context, req entity.User) (res entity.User, err error) {
	if err = req.Validate(ctx); err != nil {
		return entity.User{}, &appErr.Error{
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EValidation,
		}
	}

	err = u.UserRepo.Update(ctx, req)
	if err != nil {
		return entity.User{}, &appErr.Error{
			ErrCode: 1024,
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EConflict,
		}
	}

	return req, nil
}

func (u userService) GetByIdOrEmpty(ctx context.Context, id string) (res entity.User, err error) {
	if id == "" {
		err := errors.New("id must not be empty")
		return entity.User{}, &appErr.Error{
			ErrCode: 1024, // todo
			Cause:   err,
			Message: "Id(UserId) cannot be empty",
			Class:   appErr.EValidation,
		}
	}

	userEntity, err := u.UserRepo.FindByIdOrEmpty(ctx, id)
	if err != nil {
		return entity.User{}, &appErr.Error{
			ErrCode: 1024, // todo
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EConflict,
		}
	}

	return userEntity, nil
}

func (u userService) Purge(ctx context.Context, id string) (err error) {
	if id == "" {
		err := errors.New("id must not be empty")
		return &appErr.Error{
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EValidation,
		}
	}

	err = u.UserRepo.Purge(ctx, id)
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

func (u userService) Delete(ctx context.Context, id string) (err error) {
	if id == "" {
		err := errors.New("id must not be empty")
		return &appErr.Error{
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EValidation,
		}
	}

	err = u.UserRepo.Delete(ctx, id)
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
