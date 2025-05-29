package middleware

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thealiakbari/hichapp/pkg/common/auth"
	"github.com/thealiakbari/hichapp/pkg/common/i18next"
	appErr "github.com/thealiakbari/hichapp/pkg/common/response"
	"golang.org/x/net/context"
)

const (
	UserIdKey = "userId"
)

type middleware struct {
	umt userManagement.Transport
}

func NewMiddleware(umt userManagement.Transport) auth.Middleware {
	return &middleware{
		umt: umt,
	}
}

func (m middleware) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := ParseBearerToken(c.Request)
		if err != nil {
			appErr.HandelError(c, &appErr.Error{
				Cause:   err,
				Message: i18next.ByContext(c, "msg_authorization_fail"),
				Class:   appErr.EUnauthorized,
			})
			return
		}

		ctx := context.WithValue(c.Request.Context(), "token", token)
		res, err := m.umt.GetUserInfo(ctx)
		if err != nil {
			appErr.HandelError(c, &appErr.Error{
				Cause:   err,
				Message: i18next.ByContext(c, "msg_authorization_fail"),
				Class:   appErr.EUnauthorized,
			})
			return
		}

		c.Set(UserIdKey, res.UserId)
		c.Next()
	}
}

func (m middleware) GetUserId(c *gin.Context) (userId uuid.UUID, err error) {
	userIdStr, exist := c.Get(UserIdKey)
	if exist {
		userId, err = uuid.Parse(fmt.Sprintf("%s", userIdStr))
		if err != nil {
			return [16]byte{}, errors.New(i18next.ByContext(c, "msg_access_token_invalid"))
		}
		return userId, nil
	}

	return [16]byte{}, errors.New(i18next.ByContext(c, "msg_access_token_invalid"))
}

func (m middleware) GetUserIdFromCtx(c context.Context) (uuid.UUID, error) {
	userId, ok := c.Value(UserIdKey).(uuid.UUID)
	if ok {
		return userId, nil
	}

	return [16]byte{}, errors.New(i18next.ByContext(c, "msg_context_has_no_user_id"))
}
