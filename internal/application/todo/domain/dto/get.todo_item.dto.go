package dto

import (
	"context"

	"github.com/thealiakbari/todoapp/pkg/common/request"
)

type GetTodoItemRequest struct {
	Ids         []string `form:"ids"`
	Description []string `json:"description"`
	DueDate     []string `json:"dueDate"`

	request.Pagination `json:"-"`
}

func (g GetTodoItemRequest) Validate(ctx context.Context) error {
	return nil
}
