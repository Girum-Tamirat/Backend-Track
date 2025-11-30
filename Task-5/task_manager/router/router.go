package router

import (
	"task_manager/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(controller *controllers.TaskController) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1/tasks")
	{
		// /api/v1/tasks
		api.GET("/", controller.GetAll) 
		api.POST("/", controller.Create)
		// /api/v1/tasks/:id
		api.GET("/:id", controller.GetByID)
		api.PUT("/:id", controller.Update)
		api.DELETE("/:id", controller.Delete)
	}

	return r
}
