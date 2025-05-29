package dto

import (
	"context"

	"github.com/thealiakbari/hichapp/pkg/common/validation"
)

type UpdatePollRequest struct {
	Title string `json:"title" validate:"required"`
}

func (u UpdatePollRequest) Validate(ctx context.Context) error {
	return validation.Validate(ctx, u)
}
