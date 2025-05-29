package entity

import (
	"context"

	"github.com/thealiakbari/hichapp/pkg/common/db"
)

type User struct {
	db.UniversalModel
	Email string `gorm:"column:email;type:varchar(500);not null" validate:"required"`
}

func (u User) Validate(ctx context.Context) error {
	return nil
}
