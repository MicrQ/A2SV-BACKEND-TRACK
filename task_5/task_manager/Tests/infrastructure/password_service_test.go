package infrastructure_test

import (
	infrastructure "task_manager/Infrastructure"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordHashAndVerify(t *testing.T) {
	ps := infrastructure.NewPasswordService()
	hash, err := ps.HashPassword("secret")
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	ok := ps.VerifyPassword(hash, "secret")
	assert.True(t, ok)
	notOk := ps.VerifyPassword(hash, "wrong")
	assert.False(t, notOk)
}
