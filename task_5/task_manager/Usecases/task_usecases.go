package usecases

import (
	"time"

	domain "task_manager/Domain"
	repositories "task_manager/Repositories"
)

// TaskUsecases handles task-related business logic.
type TaskUsecases struct {
	taskRepo repositories.ITaskRepository
}

// NewTaskUsecases creates a new task usecases instance.
func NewTaskUsecases(taskRepo repositories.ITaskRepository) *TaskUsecases {
	return &TaskUsecases{taskRepo: taskRepo}
}

// GetAllTasks retrieves all tasks.
func (tu *TaskUsecases) GetAllTasks() ([]domain.Task, error) {
	return tu.taskRepo.GetAll()
}

// GetTaskByID retrieves a task by ID.
func (tu *TaskUsecases) GetTaskByID(id string) (domain.Task, error) {
	return tu.taskRepo.GetByID(id)
}

// CreateTask creates a new task after validation.
func (tu *TaskUsecases) CreateTask(title, description string, dueDate time.Time, status string) (domain.Task, error) {
	task := domain.Task{
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		Status:      status,
	}

	if err := task.Validate(); err != nil {
		return domain.Task{}, err
	}

	return tu.taskRepo.Create(task)
}

// UpdateTask updates an existing task after validation.
func (tu *TaskUsecases) UpdateTask(id, title, description string, dueDate time.Time, status string) (domain.Task, error) {
	task := domain.Task{
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		Status:      status,
	}

	if err := task.Validate(); err != nil {
		return domain.Task{}, err
	}

	return tu.taskRepo.Update(id, task)
}

// DeleteTask deletes a task by ID.
func (tu *TaskUsecases) DeleteTask(id string) error {
	return tu.taskRepo.Delete(id)
}
