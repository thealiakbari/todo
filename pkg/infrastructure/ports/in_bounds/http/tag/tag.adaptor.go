package tag

import (
	"github.com/gin-gonic/gin"
	"github.com/thealiakbari/hichapp/app/tag/service"
)

type Adaptor struct {
	service.TagHttpApp
}

func (a Adaptor) RegisterRoutes(r *gin.RouterGroup) {
	apiPoll := r.Group("/tags")

	apiPoll.POST("", a.MakeCreate())
	apiPoll.PUT("/:id", a.MakeUpdate())

	apiPoll.GET("/:id", a.MakeGetById())

	apiPoll.DELETE("/:id", a.MakeDelete())
	apiPoll.DELETE("/purge/:id", a.MakePurge())
}
