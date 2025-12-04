package middleware_test

import (
	"net/http"
	"net/http/httptest"
	domain "task_manager/Domain"
	infrastructure "task_manager/Infrastructure"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAuthMiddleware_NoAuthorization(t *testing.T) {
	jwtSvc := infrastructure.NewJWTService("secret")
	authMW := infrastructure.NewAuthMiddleware(jwtSvc)

	r := gin.New()
	r.Use(authMW.AuthRequired())
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	jwtSvc := infrastructure.NewJWTService("secret")
	authMW := infrastructure.NewAuthMiddleware(jwtSvc)

	r := gin.New()
	r.Use(authMW.AuthRequired())
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	jwtSvc := infrastructure.NewJWTService("secret")
	authMW := infrastructure.NewAuthMiddleware(jwtSvc)

	r := gin.New()
	r.Use(authMW.AuthRequired())
	r.GET("/test", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		assert.NotEmpty(t, userID)
		c.Status(200)
	})

	// Assume a valid token is generated
	token, _ := jwtSvc.GenerateToken(domain.User{ID: primitive.NewObjectID(), Username: "user", Role: "user"})
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminRequired_UserRole(t *testing.T) {
	jwtSvc := infrastructure.NewJWTService("secret")
	authMW := infrastructure.NewAuthMiddleware(jwtSvc)

	r := gin.New()
	r.Use(authMW.AuthRequired())
	r.Use(authMW.AdminRequired())
	r.GET("/admin", func(c *gin.Context) { c.Status(200) })

	token, _ := jwtSvc.GenerateToken(domain.User{ID: primitive.NewObjectID(), Username: "user", Role: "user"})
	req := httptest.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestAdminRequired_AdminRole(t *testing.T) {
	jwtSvc := infrastructure.NewJWTService("secret")
	authMW := infrastructure.NewAuthMiddleware(jwtSvc)

	r := gin.New()
	r.Use(authMW.AuthRequired())
	r.Use(authMW.AdminRequired())
	r.GET("/admin", func(c *gin.Context) { c.Status(200) })

	token, _ := jwtSvc.GenerateToken(domain.User{ID: primitive.NewObjectID(), Username: "admin", Role: "admin"})
	req := httptest.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
