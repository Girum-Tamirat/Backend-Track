package main

import (
	"log"
	"task_manager/router"
)

func main() {
	r := router.SetupRouter()

	// Run server on port 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
