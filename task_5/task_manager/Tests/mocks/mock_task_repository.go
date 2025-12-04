package mocks

import (
	domain "task_manager/Domain"

	"github.com/stretchr/testify/mock"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) GetAll() ([]domain.Task, error) {
	args := m.Called()
	if tasks, ok := args.Get(0).([]domain.Task); ok {
		return tasks, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskRepository) GetByID(id string) (domain.Task, error) {
	args := m.Called(id)
	if t, ok := args.Get(0).(domain.Task); ok {
		return t, args.Error(1)
	}
	return domain.Task{}, args.Error(1)
}

func (m *MockTaskRepository) Create(t domain.Task) (domain.Task, error) {
	args := m.Called(t)
	if created, ok := args.Get(0).(domain.Task); ok {
		return created, args.Error(1)
	}
	return domain.Task{}, args.Error(1)
}

func (m *MockTaskRepository) Update(id string, t domain.Task) (domain.Task, error) {
	args := m.Called(id, t)
	if updated, ok := args.Get(0).(domain.Task); ok {
		return updated, args.Error(1)
	}
	return domain.Task{}, args.Error(1)
}

func (m *MockTaskRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTaskRepository) Close() error {
	return nil
}
