package main

import (
	"log"
	"os"
	"task_manager/controllers"
	"task_manager/data"
	"task_manager/router"
)

func main() {
	uri := os.Getenv("MONGO_URI")
	db := os.Getenv("MONGO_DB")
	if db == "" {
		db = "taskdb"
	}

	// init services
	if err := data.InitUserService(uri, db); err != nil {
		log.Fatalf("failed to init user service: %v", err)
	}
	if err := data.InitTaskService(uri, db); err != nil {
		log.Fatalf("failed to init task service: %v", err)
	}
	defer func() {
		_ = data.GetUserService().Close()
		_ = data.GetTaskService().Close()
	}()

	ctrl := controllers.NewController(data.GetUserService(), data.GetTaskService())
	r := router.SetupRouter(ctrl)

	log.Println("Server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
