package mocks

import (
	domain "task_manager/Domain"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(username, password string) (domain.User, error) {
	args := m.Called(username, password)
	if u, ok := args.Get(0).(domain.User); ok {
		return u, args.Error(1)
	}
	return domain.User{}, args.Error(1)
}

func (m *MockUserRepository) GetByUsername(username string) (domain.User, error) {
	args := m.Called(username)
	if u, ok := args.Get(0).(domain.User); ok {
		return u, args.Error(1)
	}
	return domain.User{}, args.Error(1)
}

func (m *MockUserRepository) VerifyPassword(u domain.User, password string) bool {
	args := m.Called(u, password)
	return args.Bool(0)
}

func (m *MockUserRepository) PromoteUser(idHex string) error {
	args := m.Called(idHex)
	return args.Error(0)
}

func (m *MockUserRepository) Close() error {
	return nil
}

func (m *MockUserRepository) IsEmpty() (bool, error) {
	args := m.Called()
	return args.Bool(0), args.Error(1)
}
