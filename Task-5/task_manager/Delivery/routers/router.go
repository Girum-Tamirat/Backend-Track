package routers

import (
	"github.com/gin-gonic/gin"
	"task_manager/Delivery/controllers"
	"task_manager/Infrastructure"
)

// ctrl is passed so routes call usecases through controller
func SetupRouter(ctrl *controllers.Controller, jwtSvc Infrastructure.JWTService) *gin.Engine {
	r := gin.Default()

	// public
	r.POST("/register", ctrl.Register)
	r.POST("/login", ctrl.Login)

	// protected
	auth := r.Group("/")
	auth.Use(Infrastructure.AuthMiddleware(jwtSvc))

	// read endpoints for all authenticated users
	auth.GET("/tasks", ctrl.GetTasks)
	auth.GET("/tasks/:id", ctrl.GetTaskByID)

	// admin-only
	admin := auth.Group("/")
	admin.Use(Infrastructure.AdminOnlyMiddleware())
	admin.POST("/tasks", ctrl.CreateTask)
	admin.PUT("/tasks/:id", ctrl.UpdateTask)
	admin.DELETE("/tasks/:id", ctrl.DeleteTask)
	admin.POST("/users/:username/promote", ctrl.Promote)

	return r
}
