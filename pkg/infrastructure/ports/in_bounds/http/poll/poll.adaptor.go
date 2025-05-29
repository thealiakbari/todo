package poll

import (
	"github.com/gin-gonic/gin"
	"github.com/thealiakbari/hichapp/app/poll/service"
)

type Adaptor struct {
	service.PollHttpApp
}

func (a Adaptor) RegisterRoutes(r *gin.RouterGroup) {
	apiPoll := r.Group("/polls")

	apiPoll.POST("", a.MakeCreate())
	apiPoll.PUT("/:id", a.MakeUpdate())

	apiPoll.GET("/:id", a.MakeGetById())

	apiPoll.DELETE("/:id", a.MakeDelete())
	apiPoll.DELETE("/purge/:id", a.MakePurge())
}
