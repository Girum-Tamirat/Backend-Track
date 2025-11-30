package router

import (
	"task_manager/controllers"
	"task_manager/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(ctrl *controllers.Controller) *gin.Engine {
	r := gin.Default()

	// public auth routes
	r.POST("/register", ctrl.Register)
	r.POST("/login", ctrl.Login)

	// protected routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())

	// task routes (authenticated)
	protected.GET("/tasks", ctrl.GetTasks)
	protected.GET("/tasks/:id", ctrl.GetTaskByID)

	// admin-only task mutating routes
	admin := protected.Group("/")
	admin.Use(middleware.AdminOnly())
	admin.POST("/tasks", ctrl.CreateTask)
	admin.PUT("/tasks/:id", ctrl.UpdateTask)
	admin.DELETE("/tasks/:id", ctrl.DeleteTask)

	// admin promote endpoint
	admin.POST("/users/:username/promote", ctrl.Promote)

	return r
}
