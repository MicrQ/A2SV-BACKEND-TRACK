package usecases

import (
	"errors"

	domain "task_manager/Domain"
	repositories "task_manager/Repositories"
)

// UserUsecases handles user-related business logic.
type UserUsecases struct {
	userRepo repositories.IUserRepository
}

// NewUserUsecases creates a new user usecases instance.
func NewUserUsecases(userRepo repositories.IUserRepository) *UserUsecases {
	return &UserUsecases{userRepo: userRepo}
}

// RegisterUser registers a new user.
func (uu *UserUsecases) RegisterUser(username, password string) (domain.User, error) {
	return uu.userRepo.CreateUser(username, password)
}

// LoginUser authenticates a user and returns the user if successful.
func (uu *UserUsecases) LoginUser(username, password string) (domain.User, error) {
	user, err := uu.userRepo.GetByUsername(username)
	if err != nil {
		return domain.User{}, err
	}

	if !uu.userRepo.VerifyPassword(user, password) {
		return domain.User{}, errors.New("invalid credentials")
	}

	return user, nil
}

// PromoteUser promotes a user to admin.
func (uu *UserUsecases) PromoteUser(idHex string) error {
	return uu.userRepo.PromoteUser(idHex)
}
