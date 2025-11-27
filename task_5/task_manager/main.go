package main

import (
	"log"
	"os"

	"task_manager/Delivery/routers"
	"task_manager/Infrastructure"
	"task_manager/Repositories"
	"task_manager/Usecases"
)

func main() {
	// MongoDB configuration
	mongoURI := getEnv("MONGODB_URI", "mongodb://localhost:27017")
	dbName := getEnv("DB_NAME", "taskmanager")
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key")

	// Initialize repositories
	taskRepo, err := repositories.NewMongoTaskRepository(mongoURI, dbName, "tasks")
	if err != nil {
		log.Fatalf("Failed to connect to task repository: %v", err)
	}
	defer taskRepo.Close()

	userRepo, err := repositories.NewMongoUserRepository(mongoURI, dbName, "users")
	if err != nil {
		log.Fatalf("Failed to connect to user repository: %v", err)
	}
	defer userRepo.Close()

	// Initialize infrastructure services
	jwtService := infrastructure.NewJWTService(jwtSecret)
	authMiddleware := infrastructure.NewAuthMiddleware(jwtService)

	// Initialize usecases
	taskUsecases := usecases.NewTaskUsecases(taskRepo)
	userUsecases := usecases.NewUserUsecases(userRepo)

	// Setup router
	r := routers.SetupRouter(taskUsecases, userUsecases, jwtService, authMiddleware)

	// Run server
	log.Println("Server starting on :8080")
	r.Run(":8080")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}