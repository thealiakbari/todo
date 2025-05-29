package entity

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/thealiakbari/hichapp/pkg/common/db"
)

type Poll struct {
	db.UniversalModel
	Title string `gorm:"column:title;type:varchar(500);not null" validate:"required"`
}

func (u Poll) Validate(ctx context.Context) error {
	return nil
}

type Option struct {
	db.UniversalModel
	PollId  uuid.UUID `gorm:"column:title;type:varchar(500);not null" validate:"required"`
	Content string    `gorm:"column:content;type:text" validate:"required"`
	Counts  int32     `gorm:"column:counts;type:bigint" validate:"required"`
}

type PollTag struct {
	PollID int64
	TagID  int64
}

type Tag struct {
	ID   int64
	Name string
}

type Vote struct {
	ID        int64
	PollID    int64
	OptionID  int64
	UserID    int64
	IsSkipped bool
	CreatedAt time.Time
	UpdateAt  time.Time
	DeletedAt time.Time
}
