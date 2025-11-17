package controllers

import (
	"net/http"
	"strconv"

	"task_manager/data"
	"task_manager/models"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	service *data.TaskService
}

func NewTaskController(s *data.TaskService) *TaskController {
	return &TaskController{service: s}
}

// Register routes to a given router group
func (tc *TaskController) Register(rg *gin.RouterGroup) {
	rg.GET("/tasks", tc.ListTasks)
	rg.GET("/tasks/:id", tc.GetTask)
	rg.POST("/tasks", tc.CreateTask)
	rg.PUT("/tasks/:id", tc.UpdateTask)
	rg.DELETE("/tasks/:id", tc.DeleteTask)
}

func (tc *TaskController) ListTasks(c *gin.Context) {
	tasks := tc.service.GetAll()
	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

func (tc *TaskController) GetTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	task, err := tc.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": task})
}

func (tc *TaskController) CreateTask(c *gin.Context) {
	var input models.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	created := tc.service.Create(input)
	c.JSON(http.StatusCreated, gin.H{"data": created})
}

func (tc *TaskController) UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var input models.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, err := tc.service.Update(id, input)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": updated})
}

func (tc *TaskController) DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := tc.service.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
