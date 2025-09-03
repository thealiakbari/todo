package poll

import (
	"github.com/gin-gonic/gin"
	service "github.com/thealiakbari/todoapp/internal/application/todo"
)

type Adaptor struct {
	service.TodoItemHttpApp
}

func (a Adaptor) RegisterRoutes(r *gin.RouterGroup) {
	apiTodoItem := r.Group("/todo-items")

	apiTodoItem.POST("", a.MakeCreate())
	apiTodoItem.PUT("/:id", a.MakeUpdate())

	apiTodoItem.GET("/:id", a.MakeGetById())

	apiTodoItem.DELETE("/:id", a.MakeDelete())
	apiTodoItem.DELETE("/purge/:id", a.MakePurge())
}
