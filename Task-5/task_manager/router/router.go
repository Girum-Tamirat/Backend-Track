package router

import (
	"task_manager/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRouter initializes Gin and registers routes
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// API v1 group
	api := r.Group("/api/v1")
	taskCtrl := controllers.NewTaskController()
	taskCtrl.RegisterRoutes(api)

	return r
}
