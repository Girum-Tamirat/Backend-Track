package controllers

import (
	"net/http"

	"task_manager/data"
	"task_manager/middleware"
	"task_manager/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	UserSvc *data.UserService
	TaskSvc *data.TaskService
}

func NewController(us *data.UserService, ts *data.TaskService) *Controller {
	return &Controller{UserSvc: us, TaskSvc: ts}
}

// Register user: POST /register
func (ctr *Controller) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload", "details": err.Error()})
		return
	}
	user, err := ctr.UserSvc.CreateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"username": user.Username, "role": user.Role})
}

// Login: POST /login
func (ctr *Controller) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload", "details": err.Error()})
		return
	}
	user, err := ctr.UserSvc.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	token, err := middleware.GenerateToken(user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "username": user.Username, "role": user.Role})
}

// Promote user: POST /users/:username/promote (admin only)
func (ctr *Controller) Promote(c *gin.Context) {
	target := c.Param("username")
	if err := ctr.UserSvc.PromoteUser(target); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user promoted", "username": target})
}

// Create task: POST /tasks (admin only)
func (ctr *Controller) CreateTask(c *gin.Context) {
	var t models.Task
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload", "details": err.Error()})
		return
	}
	created, err := ctr.TaskSvc.CreateTask(t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// Get tasks: GET /tasks (authenticated users)
func (ctr *Controller) GetTasks(c *gin.Context) {
	tasks, err := ctr.TaskSvc.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tasks", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// Get task by id: GET /tasks/:id (authenticated users)
func (ctr *Controller) GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	task, err := ctr.TaskSvc.GetTaskByID(bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// Update task: PUT /tasks/:id (admin only)
func (ctr *Controller) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var payload struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		DueDate     *string `json:"due_date"`
		Status      *string `json:"status"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	
	// Build the fields to update
	updateFields := bson.M{}
	if payload.Title != nil {
		updateFields["title"] = *payload.Title
	}
	if payload.Description != nil {
		updateFields["description"] = *payload.Description
	}
	if payload.DueDate != nil {
		updateFields["due_date"] = *payload.DueDate
	}
	if payload.Status != nil {
		updateFields["status"] = *payload.Status
	}

	if len(updateFields) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	// The service layer wraps the fields with $set.
	updated, err := ctr.TaskSvc.UpdateTask(bson.M{"_id": objID}, updateFields) 
	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// Delete task: DELETE /tasks/:id (admin only)
func (ctr *Controller) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := ctr.TaskSvc.DeleteTask(bson.M{"_id": objID}); err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete", "details": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}