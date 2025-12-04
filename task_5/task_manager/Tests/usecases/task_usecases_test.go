package usecases_test

import (
	domain "task_manager/Domain"
	"task_manager/Tests/mocks"
	usecases "task_manager/Usecases"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllTasks(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	tu := usecases.NewTaskUsecases(mockRepo)

	tasks := []domain.Task{{ID: "1", Title: "Task1"}}
	mockRepo.On("GetAll").Return(tasks, nil)

	result, err := tu.GetAllTasks()
	assert.NoError(t, err)
	assert.Equal(t, tasks, result)
	mockRepo.AssertExpectations(t)
}

func TestGetTaskByID(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	tu := usecases.NewTaskUsecases(mockRepo)

	task := domain.Task{ID: "1", Title: "Task1"}
	mockRepo.On("GetByID", "1").Return(task, nil)

	result, err := tu.GetTaskByID("1")
	assert.NoError(t, err)
	assert.Equal(t, task, result)
	mockRepo.AssertExpectations(t)
}

func TestCreateTask_Valid(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	tu := usecases.NewTaskUsecases(mockRepo)

	created := domain.Task{ID: "1", Title: "New Task", Status: "pending"}
	mockRepo.On("Create", mock.AnythingOfType("domain.Task")).Return(created, nil)

	result, err := tu.CreateTask("New Task", "", time.Now(), "pending")
	assert.NoError(t, err)
	assert.Equal(t, created, result)
	mockRepo.AssertExpectations(t)
}

func TestCreateTask_InvalidStatus(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	tu := usecases.NewTaskUsecases(mockRepo)

	_, err := tu.CreateTask("Task", "", time.Now(), "invalid")
	assert.Error(t, err)
	assert.Equal(t, "invalid status", err.Error())
}

func TestUpdateTask(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	tu := usecases.NewTaskUsecases(mockRepo)

	updated := domain.Task{ID: "1", Title: "Updated"}
	mockRepo.On("Update", "1", mock.AnythingOfType("domain.Task")).Return(updated, nil)

	result, err := tu.UpdateTask("1", "Updated", "", time.Now(), "pending")
	assert.NoError(t, err)
	assert.Equal(t, updated, result)
	mockRepo.AssertExpectations(t)
}

func TestDeleteTask(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	tu := usecases.NewTaskUsecases(mockRepo)

	mockRepo.On("Delete", "1").Return(nil)

	err := tu.DeleteTask("1")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
