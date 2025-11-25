package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"task_manager/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewUserService(uri, dbName, collectionName string) (*UserService, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	coll := client.Database(dbName).Collection(collectionName)
	return &UserService{client: client, collection: coll}, nil
}

func (s *UserService) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.client.Disconnect(ctx)
}

func (s *UserService) IsEmpty() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cnt, err := s.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return false, fmt.Errorf("count error: %w", err)
	}
	return cnt == 0, nil
}

func (s *UserService) CreateUser(username, password string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check existing
	var existing models.User
	err := s.collection.FindOne(ctx, bson.M{"username": username}).Decode(&existing)
	if err == nil {
		return models.User{}, fmt.Errorf("username already exists")
	}
	if err != nil && err != mongo.ErrNoDocuments {
		return models.User{}, fmt.Errorf("failed to check username: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	role := "user"
	empty, err := s.IsEmpty()
	if err != nil {
		return models.User{}, err
	}
	if empty {
		role = "admin"
	}

	u := models.User{
		ID:           primitive.NewObjectID(),
		Username:     username,
		PasswordHash: string(hash),
		Role:         role,
		CreatedAt:    time.Now().UTC(),
	}

	if _, err := s.collection.InsertOne(ctx, u); err != nil {
		return models.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	u.PasswordHash = ""
	return u, nil
}

func (s *UserService) GetByUsername(username string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var u models.User
	if err := s.collection.FindOne(ctx, bson.M{"username": username}).Decode(&u); err != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("failed to find user: %w", err)
	}
	return u, nil
}

func (s *UserService) VerifyPassword(u models.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (s *UserService) PromoteUser(idHex string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	res, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"role": "admin"}})
	if err != nil {
		return fmt.Errorf("failed to promote: %w", err)
	}
	if res.MatchedCount == 0 {
		return ErrUserNotFound
	}
	return nil
}
