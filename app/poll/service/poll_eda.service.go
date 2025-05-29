package service

import (
	"context"
	"encoding/json"
	"errors"
	"gorm.io/gorm"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/thealiakbari/hichapp/app/poll/domain/dto"
	pollSrv "github.com/thealiakbari/hichapp/internal/poll"
	"github.com/thealiakbari/hichapp/pkg/common/db"
	"github.com/thealiakbari/hichapp/pkg/common/logger"
	appErr "github.com/thealiakbari/hichapp/pkg/common/response"
)

type consumer struct {
	db      db.DBWrapper
	log     logger.Logger
	pollSvc pollSrv.Poll
}

type Consumer interface {
	OnHichAppVoteEvent(ctx context.Context, payload message.Payload) error
}

func NewConsumer(pollSrv pollSrv.Poll, log logger.Logger, db db.DBWrapper) Consumer {
	return &consumer{
		pollSvc: pollSrv,
		log:     log,
		db:      db,
	}
}

func (c consumer) OnHichAppVoteEvent(ctx context.Context, payload message.Payload) (err error) {
	var req dto.UpdatePollRequest
	if err = json.Unmarshal(payload, &req); err != nil {
		return &appErr.Error{
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EBadArg,
		}
	}

	if err = req.Validate(ctx); err != nil {
		return &appErr.Error{
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EValidation,
		}
	}

	tx, ctx, err := db.BeginTx(ctx, c.db.DB)
	if err != nil {
		return &appErr.Error{
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EConflict,
		}
	}

	defer func() {
		if err != nil {
			if err = tx.Rollback().Error; err != nil {
				c.log.Errorf(ctx, "OnUmsUserCrudEvent err: %v", err)
			}
		}
	}()

	_, err = c.pollSvc.GetByIdOrEmpty(ctx, req.Title)
	if err != nil {
		err2 := err.(*appErr.Error)
		//  FIXME Timeout error. retry
		if !errors.Is(err2.Cause, gorm.ErrRecordNotFound) {
			return &appErr.Error{
				Cause:   err,
				Message: err.Error(),
				Class:   appErr.EDB,
			}
		}
	}

	if err = tx.Commit().Error; err != nil {
		return &appErr.Error{
			Cause:   err,
			Message: err.Error(),
			Class:   appErr.EConflict,
		}
	}

	return nil
}
