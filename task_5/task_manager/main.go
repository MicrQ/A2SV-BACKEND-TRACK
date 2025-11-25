package main

import (
	"log"
	"os"

	"task_manager/data"
	"task_manager/router"
)

func main() {
	// MongoDB configuration
	mongoURI := getEnv("MONGODB_URI", "")
	dbName := getEnv("DB_NAME", "taskmanager")
	collectionName := getEnv("COLLECTION_NAME", "tasks")

	// Create the MongoDB task service
	taskService, err := data.NewTaskService(mongoURI, dbName, collectionName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer taskService.Close()

	// Create the MongoDB user service
	userService, err := data.NewUserService(mongoURI, dbName, "users")
	if err != nil {
		log.Fatalf("Failed to create user service: %v", err)
	}
	defer userService.Close()

	// Build router (pass both services so auth routes can be registered)
	r := router.SetupRouter(taskService, userService)

	// Run server on port 8080
	log.Println("Server starting on :8080")
	r.Run(":8080")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
