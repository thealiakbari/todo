package repository

import (
	"context"
	"errors"

	"github.com/thealiakbari/hichapp/pkg/common/db"
	"github.com/thealiakbari/hichapp/pkg/common/kafka/inbox/domain/entity"
	"gorm.io/gorm"
)

type InboxMessagesRepository interface {
	Create(ctx context.Context, in entity.InboxMessage) (res entity.InboxMessage, err error)
	FindByIdAndAggregateType(ctx context.Context, id, aggregateType string) (res *entity.InboxMessage, err error)
}

type inboxMessagesRepository struct {
	db *gorm.DB
}

var instance InboxMessagesRepository = inboxMessagesRepository{}

func NewInboxMessageRepository(db *gorm.DB) inboxMessagesRepository {
	return inboxMessagesRepository{
		db: db,
	}
}

func (i inboxMessagesRepository) Create(ctx context.Context, in entity.InboxMessage) (res entity.InboxMessage, err error) {
	err = db.GormConnection(ctx, i.db).Create(&in).Error
	if err != nil {
		return entity.InboxMessage{}, err
	}

	return in, nil
}

func (i inboxMessagesRepository) FindByIdAndAggregateType(ctx context.Context, id, aggregateType string) (res *entity.InboxMessage, err error) {
	err = db.GormConnection(ctx, i.db).Model(res).Order("created_at desc").Last(&res, "id = ? AND aggregate_type = ?", id, aggregateType).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return res, nil
}
