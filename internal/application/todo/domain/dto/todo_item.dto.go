package dto

import (
	"time"

	"github.com/google/uuid"
)

type TodoItem struct {
	Id          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	DueDate     string    `json:"dueDate"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
