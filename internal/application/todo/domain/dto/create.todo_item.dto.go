package dto

import (
	"context"

	"github.com/thealiakbari/todoapp/pkg/common/validation"
)

type CreateTodoItemRequest struct {
	Description string `json:"description" validate:"required"`
	DueDate     string `json:"dueDate" validate:"required"`
}

func (c CreateTodoItemRequest) Validate(ctx context.Context) error {
	return validation.Validate(ctx, c)
}
