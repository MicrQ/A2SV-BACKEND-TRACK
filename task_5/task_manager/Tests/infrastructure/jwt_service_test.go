package infrastructure_test

import (
	d "task_manager/Domain"
	infrastructure "task_manager/Infrastructure"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestJWTGenerateAndValidate(t *testing.T) {
	jwtSvc := infrastructure.NewJWTService("test-secret")
	user := d.User{ID: primitive.NewObjectID(), Username: "u1", Role: "user"}
	token, err := jwtSvc.GenerateToken(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	claims, err := jwtSvc.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, claims["usr"])
	assert.Equal(t, user.Role, claims["role"])
}

func TestJWTInvalidToken(t *testing.T) {
	jwtSvc := infrastructure.NewJWTService("test-secret")
	_, err := jwtSvc.ValidateToken("invalid-token")
	assert.Error(t, err)
}
