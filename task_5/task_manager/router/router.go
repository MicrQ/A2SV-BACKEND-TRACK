package router

import (
	"task_manager/controllers"
	"task_manager/data"

	"github.com/gin-gonic/gin"
)

func SetupRouter(s *data.TaskService) *gin.Engine {
	r := gin.Default()

	taskController := controllers.NewTaskController(s)
	api := r.Group("/")

	taskController.Register(api)

	return r
}
