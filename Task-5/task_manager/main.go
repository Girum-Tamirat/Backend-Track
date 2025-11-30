package main

import (
	"log"
	"task_manager/controllers"
	"task_manager/data"
	"task_manager/router"
)

func main() {

	service, err := data.NewTaskService()
	if err != nil {
		// Use log.Fatalf to exit if DB connection fails
		log.Fatalf("MongoDB connection failed: %v", err)
	}

	controller := controllers.NewTaskController(service)

	app := router.SetupRoutes(controller)

	log.Println("Server running on http://localhost:8080...")
	// The application will run on port 8080
	if err := app.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
