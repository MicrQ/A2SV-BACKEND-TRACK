package controllers

import (
	"net/http"

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
	tasks, err := tc.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve tasks"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

func (tc *TaskController) GetTask(c *gin.Context) {
	idStr := c.Param("id")
	task, err := tc.service.GetByID(idStr)
	if err != nil {
		if err == data.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve task"})
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
	created, err := tc.service.Create(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": created})
}

func (tc *TaskController) UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	var input models.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, err := tc.service.Update(idStr, input)
	if err != nil {
		if err == data.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": updated})
}

func (tc *TaskController) DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	if err := tc.service.Delete(idStr); err != nil {
		if err == data.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task"})
		return
	}
	c.Status(http.StatusNoContent)
}
