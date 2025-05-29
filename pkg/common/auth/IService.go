package auth

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Middleware interface {
	Authorize() gin.HandlerFunc
	GetUserId(c *gin.Context) (uuid.UUID, error)
	GetUserIdFromCtx(c context.Context) (uuid.UUID, error)
}
