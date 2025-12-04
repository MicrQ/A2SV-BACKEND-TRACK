package repositories_integration_test

import (
	"context"
	"os"
	"testing"
	"time"

	repositories "task_manager/Repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepositoryIntegrationSuite struct {
	suite.Suite
	repo       repositories.IUserRepository
	client     *mongo.Client
	collection *mongo.Collection
}

func (s *UserRepositoryIntegrationSuite) SetupSuite() {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	repo, err := repositories.NewMongoUserRepository(mongoURI, "taskmanager_test", "users_test")
	s.Require().NoError(err)
	s.repo = repo

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	s.Require().NoError(err)
	s.client = client
	s.collection = client.Database("taskmanager_test").Collection("users_test")
}

func (s *UserRepositoryIntegrationSuite) TearDownSuite() {
	if s.repo != nil {
		s.repo.Close()
	}
	if s.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.client.Disconnect(ctx)
	}
}

func (s *UserRepositoryIntegrationSuite) SetupTest() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.collection.DeleteMany(ctx, bson.M{})
}

func (s *UserRepositoryIntegrationSuite) TestCreateUser_FirstUserIsAdmin() {
	user, err := s.repo.CreateUser("firstuser", "password123")
	assert.NoError(s.T(), err)
	assert.NotEqual(s.T(), primitive.NilObjectID, user.ID)
	assert.Equal(s.T(), "firstuser", user.Username)
	assert.Equal(s.T(), "admin", user.Role)
}

func (s *UserRepositoryIntegrationSuite) TestCreateUser_SubsequentUserIsRegular() {
	_, err := s.repo.CreateUser("admin", "password123")
	assert.NoError(s.T(), err)

	user, err := s.repo.CreateUser("regularuser", "password123")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "user", user.Role)
}

func (s *UserRepositoryIntegrationSuite) TestCreateUser_DuplicateUsername() {
	_, err := s.repo.CreateUser("duplicate", "password123")
	assert.NoError(s.T(), err)

	_, err = s.repo.CreateUser("duplicate", "password456")
	assert.Error(s.T(), err)
}

func (s *UserRepositoryIntegrationSuite) TestGetByUsername() {
	created, err := s.repo.CreateUser("findme", "password123")
	assert.NoError(s.T(), err)

	found, err := s.repo.GetByUsername("findme")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), created.Username, found.Username)
}

func (s *UserRepositoryIntegrationSuite) TestGetByUsername_NotFound() {
	_, err := s.repo.GetByUsername("nonexistent")
	assert.Error(s.T(), err)
}

func (s *UserRepositoryIntegrationSuite) TestVerifyPassword() {
	_, err := s.repo.CreateUser("verifyuser", "correctpassword")
	assert.NoError(s.T(), err)

	user, err := s.repo.GetByUsername("verifyuser")
	assert.NoError(s.T(), err)

	assert.True(s.T(), s.repo.VerifyPassword(user, "correctpassword"))
	assert.False(s.T(), s.repo.VerifyPassword(user, "wrongpassword"))
}

func (s *UserRepositoryIntegrationSuite) TestPromoteUser() {
	_, err := s.repo.CreateUser("admin", "password123")
	assert.NoError(s.T(), err)

	user, err := s.repo.CreateUser("topromote", "password123")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "user", user.Role)

	err = s.repo.PromoteUser(user.ID.Hex())
	assert.NoError(s.T(), err)

	promoted, err := s.repo.GetByUsername("topromote")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "admin", promoted.Role)
}

func TestUserRepositoryIntegrationSuite(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION") == "true" {
		t.Skip("Skipping integration tests")
	}
	suite.Run(t, new(UserRepositoryIntegrationSuite))
}
