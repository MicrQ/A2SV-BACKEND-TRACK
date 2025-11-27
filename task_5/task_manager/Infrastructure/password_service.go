package infrastructure

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordService handles password hashing and verification.
type PasswordService struct{}

// NewPasswordService creates a new password service.
func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

// HashPassword hashes the given password.
func (p *PasswordService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword verifies the password against the hash.
func (p *PasswordService) VerifyPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
