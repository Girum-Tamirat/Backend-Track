package controllers

import (
	"net/http"
	"strconv"
	"task_manager/data"
	"task_manager/models"

	"github.com/gin-gonic/gin"
)

// TaskController holds dependencies for task handlers
type TaskController struct {
	service *data.TaskService 
}

// NewTaskController returns a new controller wired to the service singleton
func NewTaskController() *TaskController {
	return &TaskController{
		service: data.GetService(),
	}
}

// RegisterRoutes registers the controller routes on the provided router group
func (tc *TaskController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/tasks", tc.GetAll)
	rg.GET("/tasks/:id", tc.GetByID)
	rg.POST("/tasks", tc.Create)
	rg.PUT("/tasks/:id", tc.Update)
	rg.DELETE("/tasks/:id", tc.Delete)
}

// GetAll - GET /tasks
func (tc *TaskController) GetAll(c *gin.Context) {
	tasks := tc.service.GetAll()
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// GetByID - GET /tasks/:id
func (tc *TaskController) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	task, err := tc.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// Create - POST /tasks
func (tc *TaskController) Create(c *gin.Context) {
	var input models.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	created := tc.service.Create(input)
	c.JSON(http.StatusCreated, created)
}

// Update - PUT /tasks/:id
func (tc *TaskController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	var input models.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	updated, err := tc.service.Update(id, input)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// Delete - DELETE /tasks/:id
func (tc *TaskController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	if err := tc.service.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
