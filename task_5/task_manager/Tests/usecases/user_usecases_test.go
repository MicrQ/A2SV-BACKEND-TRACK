package usecases_test

import (
	domain "task_manager/Domain"
	"task_manager/Tests/mocks"
	usecases "task_manager/Usecases"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRegisterUser_FirstUserAdmin(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	uu := usecases.NewUserUsecases(mockRepo)

	mockRepo.On("CreateUser", "admin", "pass").Return(domain.User{ID: primitive.NewObjectID(), Username: "admin", Role: "admin"}, nil)

	user, err := uu.RegisterUser("admin", "pass")
	assert.NoError(t, err)
	assert.Equal(t, "admin", user.Role)
	mockRepo.AssertExpectations(t)
}

func TestRegisterUser_SubsequentUser(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	uu := usecases.NewUserUsecases(mockRepo)

	mockRepo.On("CreateUser", "user", "pass").Return(domain.User{ID: primitive.NewObjectID(), Username: "user", Role: "user"}, nil)

	user, err := uu.RegisterUser("user", "pass")
	assert.NoError(t, err)
	assert.Equal(t, "user", user.Role)
	mockRepo.AssertExpectations(t)
}

func TestLoginUser_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	uu := usecases.NewUserUsecases(mockRepo)

	user := domain.User{ID: primitive.NewObjectID(), Username: "user", PasswordHash: "hash", Role: "user"}
	mockRepo.On("GetByUsername", "user").Return(user, nil)
	mockRepo.On("VerifyPassword", user, "pass").Return(true)

	loggedIn, err := uu.LoginUser("user", "pass")
	assert.NoError(t, err)
	assert.Equal(t, user.Username, loggedIn.Username)
	mockRepo.AssertExpectations(t)
}

func TestLoginUser_InvalidCredentials(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	uu := usecases.NewUserUsecases(mockRepo)

	user := domain.User{ID: primitive.NewObjectID(), Username: "user", PasswordHash: "hash", Role: "user"}
	mockRepo.On("GetByUsername", "user").Return(user, nil)
	mockRepo.On("VerifyPassword", user, "pass").Return(false)

	_, err := uu.LoginUser("user", "pass")
	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestPromoteUser(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	uu := usecases.NewUserUsecases(mockRepo)

	mockRepo.On("PromoteUser", "id").Return(nil)

	err := uu.PromoteUser("id")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
