package dto

import (
	"context"

	"github.com/thealiakbari/hichapp/pkg/common/validation"
)

type CreateTagRequest struct {
	Email string `json:"email" validate:"omitempty,email"`
}

func (c CreateTagRequest) Validate(ctx context.Context) error {
	return validation.Validate(ctx, c)
}
