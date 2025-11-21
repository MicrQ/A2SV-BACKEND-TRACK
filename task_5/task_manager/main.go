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

	// Build router
	r := router.SetupRouter(taskService)

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
