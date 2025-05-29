package dto

import (
	"context"

	"github.com/thealiakbari/hichapp/pkg/common/request"
)

type GetPollRequest struct {
	Ids   []string `form:"ids"`
	Title []string `form:"title"`

	request.Pagination `json:"-"`
}

func (g GetPollRequest) Validate(ctx context.Context) error {
	return nil
}
