package controllers

import (
	"net/http"
	"os"
	"time"

	"task_manager/data"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserController struct {
	service *data.UserService
	jwtKey  []byte
}

func NewUserController(s *data.UserService) *UserController {
	secret := os.Getenv("JWT_SECRET")
	return &UserController{service: s, jwtKey: []byte(secret)}
}

// RegisterRoutes user/auth routes
func (uc *UserController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/register", uc.Register)
	rg.POST("/login", uc.Login)
	rg.POST("/promote/:id", uc.Promote)
}

// Register creates a new user
func (uc *UserController) Register(c *gin.Context) {
	var in struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := uc.service.CreateUser(in.Username, in.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": u})
}

// Login authenticates and returns a JWT token
func (uc *UserController) Login(c *gin.Context) {
	var in struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := uc.service.GetByUsername(in.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if !uc.service.VerifyPassword(u, in.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if len(uc.jwtKey) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT secret not configured"})
		return
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  u.ID.Hex(),
		"usr":  u.Username,
		"role": u.Role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(uc.jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// Promote promotes a user to admin (caller should be protected by middleware)
func (uc *UserController) Promote(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}
	if err := uc.service.PromoteUser(id); err != nil {
		if err == data.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to promote user"})
		return
	}
	c.Status(http.StatusNoContent)
}
