package router

import (
	"task_manager/controllers"
	"task_manager/data"

	"github.com/gin-gonic/gin"
)

// SetupRouter initializes routes. Accept both task and user services so auth
// routes and middleware can be registered.
func SetupRouter(s *data.TaskService, u *data.UserService) *gin.Engine {
	r := gin.Default()

	taskController := controllers.NewTaskController(s)
	userController := controllers.NewUserController(u)

	api := r.Group("/")

	// Register task routes
	taskController.Register(api)

	// Register auth/user routes
	userController.RegisterRoutes(api)

	return r
}
