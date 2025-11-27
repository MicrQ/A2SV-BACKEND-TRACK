package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	domain "task_manager/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrNotFound = errors.New("task not found")
)

// ITaskRepository defines the interface for task data access.
type ITaskRepository interface {
	GetAll() ([]domain.Task, error)
	GetByID(id string) (domain.Task, error)
	Create(t domain.Task) (domain.Task, error)
	Update(id string, t domain.Task) (domain.Task, error)
	Delete(id string) error
	Close() error
}

// MongoTaskRepository implements ITaskRepository using MongoDB.
type MongoTaskRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoTaskRepository(uri string, dbName string, collectionName string) (ITaskRepository, error) {
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

	return &MongoTaskRepository{
		client:     client,
		collection: collection,
	}, nil
}

func (r *MongoTaskRepository) GetAll() ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find tasks: %v", err)
	}
	defer cursor.Close(ctx)

	var tasks []domain.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, fmt.Errorf("failed to decode tasks: %v", err)
	}

	return tasks, nil
}

func (r *MongoTaskRepository) GetByID(id string) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var task domain.Task
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Task{}, ErrNotFound
		}
		return domain.Task{}, fmt.Errorf("failed to find task: %v", err)
	}

	return task, nil
}

func (r *MongoTaskRepository) Create(t domain.Task) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Generate a new ID if not provided
	if t.ID == "" {
		t.ID = primitive.NewObjectID().Hex()
	}

	_, err := r.collection.InsertOne(ctx, t)
	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to insert task: %v", err)
	}

	return t, nil
}

func (r *MongoTaskRepository) Update(id string, t domain.Task) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.ID = id
	update := bson.M{"$set": t}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to update task: %v", err)
	}

	if result.MatchedCount == 0 {
		return domain.Task{}, ErrNotFound
	}

	return t, nil
}

func (r *MongoTaskRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete task: %v", err)
	}

	if result.DeletedCount == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *MongoTaskRepository) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.client.Disconnect(ctx)
}
