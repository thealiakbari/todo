package dto

import (
	"context"

	"github.com/thealiakbari/hichapp/pkg/common/validation"
)

type UpdateTagRequest struct {
	Email string `json:"email" validate:"omitempty,email"`
}

func (u UpdateTagRequest) Validate(ctx context.Context) error {
	return validation.Validate(ctx, u)
}
