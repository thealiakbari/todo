package dto

import (
	"context"

	"github.com/thealiakbari/hichapp/pkg/common/validation"
)

type CreateUserRequest struct {
	Email string `json:"email" validate:"omitempty,email"`
}

func (c CreateUserRequest) Validate(ctx context.Context) error {
	return validation.Validate(ctx, c)
}
