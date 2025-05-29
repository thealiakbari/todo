package user

import (
	"github.com/gin-gonic/gin"
	"github.com/thealiakbari/hichapp/app/user/service"
)

type Adaptor struct {
	service.UserHttpApp
}

func (a Adaptor) RegisterRoutes(r *gin.RouterGroup) {
	apiUser := r.Group("/users")

	apiUser.POST("", a.MakeCreate())
	apiUser.PUT("/:id", a.MakeUpdate())

	apiUser.GET("/:id", a.MakeGetById())

	apiUser.DELETE("/:id", a.MakeDelete())
	apiUser.DELETE("/purge/:id", a.MakePurge())
}
