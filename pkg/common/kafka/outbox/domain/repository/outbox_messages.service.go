package repository

import (
	"context"

	"github.com/thealiakbari/hichapp/pkg/common/db"
	"github.com/thealiakbari/hichapp/pkg/common/kafka/outbox/domain/entity"
	"gorm.io/gorm"
)

type OutboxMessagesRepository interface {
	Create(ctx context.Context, in entity.OutboxMessage) (res entity.OutboxMessage, err error)
}

type outboxMessagesRepository struct {
	db *gorm.DB
}

var instance OutboxMessagesRepository = outboxMessagesRepository{}

func NewInboxMessageRepository(db *gorm.DB) outboxMessagesRepository {
	return outboxMessagesRepository{
		db: db,
	}
}

func (i outboxMessagesRepository) Create(ctx context.Context, in entity.OutboxMessage) (res entity.OutboxMessage, err error) {
	err = db.GormConnection(ctx, i.db).Create(&in).Error
	if err != nil {
		return entity.OutboxMessage{}, err
	}

	return in, nil
}
