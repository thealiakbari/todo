package dto

import (
	"context"

	"github.com/thealiakbari/todoapp/pkg/common/validation"
)

type UpdateTodoItemRequest struct {
	Description string `json:"description" validate:"required"`
	DueDate     string `json:"dueDate" validate:"required"`
}

func (u UpdateTodoItemRequest) Validate(ctx context.Context) error {
	return validation.Validate(ctx, u)
}
