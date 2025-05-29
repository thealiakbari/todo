package entity

import (
	"context"

	"github.com/google/uuid"
	"github.com/thealiakbari/hichapp/pkg/common/db"
)

type Tag struct {
	db.UniversalModel

	Name string `gorm:"column:name;type:varchar(500);not null" validate:"required"`
}

func (u Tag) Validate(ctx context.Context) error {
	return nil
}

type Option struct {
	db.UniversalModel

	TagId   uuid.UUID `gorm:"column:tag_id;type:uuid;not null" validate:"required"`
	Content string    `gorm:"column:content;type:text" validate:"required"`
	Counts  int32     `gorm:"column:counts;type:bigint" validate:"required"`
}

func (u Option) Validate(ctx context.Context) error {
	return nil
}

type PollTag struct {
	db.UniversalModel
	TagId  uuid.UUID `gorm:"column:tag_id;type:uuid;not null" validate:"required"`
	PollId uuid.UUID `gorm:"column:poll_id;type:uuid;not null" validate:"required"`
}

func (u PollTag) Validate(ctx context.Context) error {
	return nil
}

type Vote struct {
	db.UniversalModel

	UserId   uuid.UUID `gorm:"column:user_id;type:uuid;not null" validate:"required"`
	TagId    uuid.UUID `gorm:"column:tag_id;type:uuid;not null" validate:"required"`
	OptionId uuid.UUID `gorm:"column:option_id;type:uuid;not null" validate:"required"`
}

func (u Vote) Validate(ctx context.Context) error {
	return nil
}
