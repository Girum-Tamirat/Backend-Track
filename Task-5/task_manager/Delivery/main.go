package main

import (
	"log"
	"os"

	"task_manager/Delivery/controllers"
	"task_manager/Delivery/routers"
	"task_manager/Repositories/mongoimpl"
	"task_manager/Usecases"
	"task_manager/Infrastructure"
)

func main() {
	// Config via env
	mongoURI := os.Getenv("MONGO_URI")
	mongoDB := os.Getenv("MONGO_DB")
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "change_this_secret"
	}

	// connect repositories (Mongo)
	mongoClient, err := mongoimpl.NewMongoClient(mongoURI, mongoDB)
	if err != nil {
		log.Fatalf("failed to connect mongo: %v", err)
	}
	defer mongoClient.Close()

	taskRepo := mongoimpl.NewTaskRepository(mongoClient)
	userRepo := mongoimpl.NewUserRepository(mongoClient)

	// usecases
	userUC := Usecases.NewUserUsecase(userRepo)
	taskUC := Usecases.NewTaskUsecase(taskRepo)

	// infrastructure (jwt service)
	infraJwt := Infrastructure.NewJWTService(jwtSecret)

	// controller
	ctrl := controllers.NewController(userUC, taskUC, infraJwt)

	// router
	r := routers.SetupRouter(ctrl, infraJwt)

	log.Println("Server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
