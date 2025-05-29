package dto

import (
	"time"

	"github.com/google/uuid"
)

type Poll struct {
	Id        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
