package dto

import (
	"context"

	"github.com/thealiakbari/hichapp/pkg/common/validation"
)

type UpdateUserRequest struct {
	Email string `json:"email" validate:"omitempty,email"`
}

func (u UpdateUserRequest) Validate(ctx context.Context) error {
	return validation.Validate(ctx, u)
}
