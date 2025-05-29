package dto

import (
	"context"

	"github.com/thealiakbari/hichapp/pkg/common/validation"
)

type CreatePollRequest struct {
	Title string `json:"title" validate:"required"`
}

func (c CreatePollRequest) Validate(ctx context.Context) error {
	return validation.Validate(ctx, c)
}
