package controllers

import (
	"net/http"
	"time"

	infrastructure "task_manager/Infrastructure"
	usecases "task_manager/Usecases"

	"github.com/gin-gonic/gin"
)

// Controller handles HTTP requests and responses.
type Controller struct {
	taskUsecases *usecases.TaskUsecases
	userUsecases *usecases.UserUsecases
	jwtService   *infrastructure.JWTService
}

// NewController creates a new controller.
func NewController(taskUsecases *usecases.TaskUsecases, userUsecases *usecases.UserUsecases, jwtService *infrastructure.JWTService) *Controller {
	return &Controller{
		taskUsecases: taskUsecases,
		userUsecases: userUsecases,
		jwtService:   jwtService,
	}
}

// Task Handlers

// ListTasks handles GET /tasks
func (c *Controller) ListTasks(ctx *gin.Context) {
	tasks, err := c.taskUsecases.GetAllTasks()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve tasks"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": tasks})
}

// GetTask handles GET /tasks/:id
func (c *Controller) GetTask(ctx *gin.Context) {
	id := ctx.Param("id")
	task, err := c.taskUsecases.GetTaskByID(id)
	if err != nil {
		if err.Error() == "task not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve task"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": task})
}

// CreateTask handles POST /tasks
func (c *Controller) CreateTask(ctx *gin.Context) {
	var input struct {
		Title       string    `json:"title" binding:"required"`
		Description string    `json:"description"`
		DueDate     time.Time `json:"due_date"`
		Status      string    `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := c.taskUsecases.CreateTask(input.Title, input.Description, input.DueDate, input.Status)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"data": task})
}

// UpdateTask handles PUT /tasks/:id
func (c *Controller) UpdateTask(ctx *gin.Context) {
	id := ctx.Param("id")
	var input struct {
		Title       string    `json:"title" binding:"required"`
		Description string    `json:"description"`
		DueDate     time.Time `json:"due_date"`
		Status      string    `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := c.taskUsecases.UpdateTask(id, input.Title, input.Description, input.DueDate, input.Status)
	if err != nil {
		if err.Error() == "task not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": task})
}

// DeleteTask handles DELETE /tasks/:id
func (c *Controller) DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.taskUsecases.DeleteTask(id); err != nil {
		if err.Error() == "task not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task"})
		return
	}
	ctx.Status(http.StatusNoContent)
}

// User/Auth Handlers

// Register handles POST /register
func (c *Controller) Register(ctx *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userUsecases.RegisterUser(input.Username, input.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"data": user})
}

// Login handles POST /login
func (c *Controller) Login(ctx *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userUsecases.LoginUser(input.Username, input.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := c.jwtService.GenerateToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

// Promote handles POST /promote/:id
func (c *Controller) Promote(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.userUsecases.PromoteUser(id); err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to promote user"})
		return
	}
	ctx.Status(http.StatusNoContent)
}
