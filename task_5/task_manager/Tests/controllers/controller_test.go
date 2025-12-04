package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"task_manager/Delivery/controllers"
	domain "task_manager/Domain"
	"task_manager/Tests/mocks"
	usecases "task_manager/Usecases"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestController_ListTasks(t *testing.T) {
	mockTaskRepo := new(mocks.MockTaskRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	taskUsecases := usecases.NewTaskUsecases(mockTaskRepo)
	userUsecases := usecases.NewUserUsecases(mockUserRepo)
	ctrl := controllers.NewController(taskUsecases, userUsecases, nil) // jwt not needed for this test

	tasks := []domain.Task{{ID: "1", Title: "Task1"}}
	mockTaskRepo.On("GetAll").Return(tasks, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	ctrl.ListTasks(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string][]domain.Task
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, tasks, response["data"])
	mockTaskRepo.AssertExpectations(t)
}

func TestController_GetTask(t *testing.T) {
	mockTaskRepo := new(mocks.MockTaskRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	taskUsecases := usecases.NewTaskUsecases(mockTaskRepo)
	userUsecases := usecases.NewUserUsecases(mockUserRepo)
	ctrl := controllers.NewController(taskUsecases, userUsecases, nil)

	task := domain.Task{ID: "1", Title: "Task1"}
	mockTaskRepo.On("GetByID", "1").Return(task, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	ctrl.GetTask(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]domain.Task
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, task, response["data"])
	mockTaskRepo.AssertExpectations(t)
}

func TestController_CreateTask(t *testing.T) {
	mockTaskRepo := new(mocks.MockTaskRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	taskUsecases := usecases.NewTaskUsecases(mockTaskRepo)
	userUsecases := usecases.NewUserUsecases(mockUserRepo)
	ctrl := controllers.NewController(taskUsecases, userUsecases, nil)

	created := domain.Task{ID: "1", Title: "New Task", Status: "pending"}
	mockTaskRepo.On("Create", mock.AnythingOfType("domain.Task")).Return(created, nil)

	body := map[string]interface{}{
		"title":  "New Task",
		"status": "pending",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	ctrl.CreateTask(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]domain.Task
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, created, response["data"])
	mockTaskRepo.AssertExpectations(t)
}

func TestController_Register(t *testing.T) {
	mockTaskRepo := new(mocks.MockTaskRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	taskUsecases := usecases.NewTaskUsecases(mockTaskRepo)
	userUsecases := usecases.NewUserUsecases(mockUserRepo)
	ctrl := controllers.NewController(taskUsecases, userUsecases, nil)

	user := domain.User{Username: "user", Role: "user"}
	mockUserRepo.On("CreateUser", "user", "pass").Return(user, nil)

	body := map[string]string{
		"username": "user",
		"password": "pass",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	ctrl.Register(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]domain.User
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, user, response["data"])
	mockUserRepo.AssertExpectations(t)
}
