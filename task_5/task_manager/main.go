package main

import (
	"task_manager/data"
	"task_manager/router"
)

func main() {
	// Create the in-memory task service
	taskService := data.NewTaskService()

	// Build router
	r := router.SetupRouter(taskService)

	// Run server on port 8080
	r.Run()
}