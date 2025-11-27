package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	domain "task_manager/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

// IUserRepository defines the interface for user data access.
type IUserRepository interface {
	CreateUser(username, password string) (domain.User, error)
	GetByUsername(username string) (domain.User, error)
	VerifyPassword(u domain.User, password string) bool
	PromoteUser(idHex string) error
	Close() error
	IsEmpty() (bool, error)
}

// MongoUserRepository implements IUserRepository using MongoDB.
type MongoUserRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoUserRepository(uri, dbName, collectionName string) (IUserRepository, error) {
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
	return &MongoUserRepository{client: client, collection: coll}, nil
}

func (r *MongoUserRepository) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.client.Disconnect(ctx)
}

func (r *MongoUserRepository) IsEmpty() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cnt, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return false, fmt.Errorf("count error: %w", err)
	}
	return cnt == 0, nil
}

func (r *MongoUserRepository) CreateUser(username, password string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check existing
	var existing domain.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&existing)
	if err == nil {
		return domain.User{}, fmt.Errorf("username already exists")
	}
	if err != nil && err != mongo.ErrNoDocuments {
		return domain.User{}, fmt.Errorf("failed to check username: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	role := "user"
	empty, err := r.IsEmpty()
	if err != nil {
		return domain.User{}, err
	}
	if empty {
		role = "admin"
	}

	u := domain.User{
		ID:           primitive.NewObjectID(),
		Username:     username,
		PasswordHash: string(hash),
		Role:         role,
		CreatedAt:    time.Now().UTC(),
	}

	if _, err := r.collection.InsertOne(ctx, u); err != nil {
		return domain.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	u.PasswordHash = ""
	return u, nil
}

func (r *MongoUserRepository) GetByUsername(username string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var u domain.User
	if err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&u); err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.User{}, ErrUserNotFound
		}
		return domain.User{}, fmt.Errorf("failed to find user: %w", err)
	}
	return u, nil
}

func (r *MongoUserRepository) VerifyPassword(u domain.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (r *MongoUserRepository) PromoteUser(idHex string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	res, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"role": "admin"}})
	if err != nil {
		return fmt.Errorf("failed to promote: %w", err)
	}
	if res.MatchedCount == 0 {
		return ErrUserNotFound
	}
	return nil
}
