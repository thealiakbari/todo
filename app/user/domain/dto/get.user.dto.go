package dto

import (
	"context"

	"github.com/thealiakbari/hichapp/pkg/common/request"
)

type GetUserRequest struct {
	Ids    []string `form:"ids"`
	Emails []string `form:"emails"`

	request.Pagination `json:"-"`
}

func (g GetUserRequest) Validate(ctx context.Context) error {
	return nil
}
