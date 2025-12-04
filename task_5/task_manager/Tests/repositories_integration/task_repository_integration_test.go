package repositories_integration_test

import (
	"context"
	"os"
	"testing"
	"time"

	domain "task_manager/Domain"
	repositories "task_manager/Repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskRepositoryIntegrationSuite struct {
	suite.Suite
	repo       repositories.ITaskRepository
	client     *mongo.Client
	collection *mongo.Collection
}

func (s *TaskRepositoryIntegrationSuite) SetupSuite() {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	repo, err := repositories.NewMongoTaskRepository(mongoURI, "taskmanager_test", "tasks_test")
	s.Require().NoError(err)
	s.repo = repo

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	s.Require().NoError(err)
	s.client = client
	s.collection = client.Database("taskmanager_test").Collection("tasks_test")
}

func (s *TaskRepositoryIntegrationSuite) TearDownSuite() {
	if s.repo != nil {
		s.repo.Close()
	}
	if s.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.client.Disconnect(ctx)
	}
}

func (s *TaskRepositoryIntegrationSuite) SetupTest() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.collection.DeleteMany(ctx, bson.M{})
}

func (s *TaskRepositoryIntegrationSuite) TestCreateAndGetTask() {
	task := domain.Task{
		Title:       "Integration Test Task",
		Description: "Testing MongoDB integration",
		Status:      "pending",
		DueDate:     time.Now().Add(24 * time.Hour),
	}

	created, err := s.repo.Create(task)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), created.ID)
	assert.Equal(s.T(), task.Title, created.Title)

	fetched, err := s.repo.GetByID(created.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), created.ID, fetched.ID)
	assert.Equal(s.T(), task.Title, fetched.Title)
}

func (s *TaskRepositoryIntegrationSuite) TestGetAllTasks() {
	task1 := domain.Task{Title: "Task 1", Status: "pending"}
	task2 := domain.Task{Title: "Task 2", Status: "completed"}

	_, err := s.repo.Create(task1)
	assert.NoError(s.T(), err)
	_, err = s.repo.Create(task2)
	assert.NoError(s.T(), err)

	tasks, err := s.repo.GetAll()
	assert.NoError(s.T(), err)
	assert.Len(s.T(), tasks, 2)
}

func (s *TaskRepositoryIntegrationSuite) TestUpdateTask() {
	task := domain.Task{Title: "Original", Status: "pending"}
	created, err := s.repo.Create(task)
	assert.NoError(s.T(), err)

	updated, err := s.repo.Update(created.ID, domain.Task{
		Title:  "Updated",
		Status: "completed",
	})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "Updated", updated.Title)

	fetched, err := s.repo.GetByID(created.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "Updated", fetched.Title)
	assert.Equal(s.T(), "completed", fetched.Status)
}

func (s *TaskRepositoryIntegrationSuite) TestDeleteTask() {
	task := domain.Task{Title: "To Delete", Status: "pending"}
	created, err := s.repo.Create(task)
	assert.NoError(s.T(), err)

	err = s.repo.Delete(created.ID)
	assert.NoError(s.T(), err)

	_, err = s.repo.GetByID(created.ID)
	assert.Error(s.T(), err)
}

func (s *TaskRepositoryIntegrationSuite) TestGetByID_NotFound() {
	_, err := s.repo.GetByID("nonexistent-id")
	assert.Error(s.T(), err)
}

func TestTaskRepositoryIntegrationSuite(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION") == "true" {
		t.Skip("Skipping integration tests")
	}
	suite.Run(t, new(TaskRepositoryIntegrationSuite))
}
