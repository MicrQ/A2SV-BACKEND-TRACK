package routers

import (
	"task_manager/Delivery/controllers"
	infrastructure "task_manager/Infrastructure"
	usecases "task_manager/Usecases"

	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the Gin router with routes and middleware.
func SetupRouter(taskUsecases *usecases.TaskUsecases, userUsecases *usecases.UserUsecases, jwtService *infrastructure.JWTService, authMiddleware *infrastructure.AuthMiddleware) *gin.Engine {
	r := gin.Default()

	ctrl := controllers.NewController(taskUsecases, userUsecases, jwtService)

	// Public routes
	r.POST("/register", ctrl.Register)
	r.POST("/login", ctrl.Login)

	// Protected routes
	protected := r.Group("/")
	protected.Use(authMiddleware.AuthRequired())
	{
		protected.GET("/tasks", ctrl.ListTasks)
		protected.GET("/tasks/:id", ctrl.GetTask)
		protected.POST("/tasks", ctrl.CreateTask)
		protected.PUT("/tasks/:id", ctrl.UpdateTask)
		protected.DELETE("/tasks/:id", authMiddleware.AdminRequired(), ctrl.DeleteTask)
		protected.POST("/promote/:id", authMiddleware.AdminRequired(), ctrl.Promote)
	}

	return r
}
