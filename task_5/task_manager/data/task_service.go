package data

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"task_manager/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrNotFound = errors.New("task not found")
)

type TaskService struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewTaskService(uri string, dbName string, collectionName string) (*TaskService, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB")

	collection := client.Database(dbName).Collection(collectionName)

	return &TaskService{
		client:     client,
		collection: collection,
	}, nil
}

func (s *TaskService) GetAll() ([]models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find tasks: %v", err)
	}
	defer cursor.Close(ctx)

	var tasks []models.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, fmt.Errorf("failed to decode tasks: %v", err)
	}

	return tasks, nil
}

func (s *TaskService) GetByID(id string) (models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var task models.Task
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Task{}, ErrNotFound
		}
		return models.Task{}, fmt.Errorf("failed to find task: %v", err)
	}

	return task, nil
}

func (s *TaskService) Create(t models.Task) (models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Generate a new ID if not provided
	if t.ID == "" {
		t.ID = primitive.NewObjectID().Hex()
	}

	_, err := s.collection.InsertOne(ctx, t)
	if err != nil {
		return models.Task{}, fmt.Errorf("failed to insert task: %v", err)
	}

	return t, nil
}

func (s *TaskService) Update(id string, t models.Task) (models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.ID = id
	update := bson.M{"$set": t}

	result, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return models.Task{}, fmt.Errorf("failed to update task: %v", err)
	}

	if result.MatchedCount == 0 {
		return models.Task{}, ErrNotFound
	}

	return t, nil
}

func (s *TaskService) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete task: %v", err)
	}

	if result.DeletedCount == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *TaskService) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.client.Disconnect(ctx)
}

