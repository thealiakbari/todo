package entity

import (
	"context"
	"github.com/thealiakbari/todoapp/pkg/common/db"
)

type TodoItem struct {
	db.UniversalModel
	Description string `gorm:"column:description;type:text;not null" validate:"required"`
	DueDate     string `gorm:"column:due_date;type:timestamp;not null" validate:"required"`
}

func (u TodoItem) Validate(ctx context.Context) error {
	return nil
}
