package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"task_manager/Domain"
	"task_manager/Infrastructure"
	"task_manager/Usecases"
)

type Controller struct {
	userUC   Usecases.UserUsecase
	taskUC   Usecases.TaskUsecase
	jwtSvc   Infrastructure.JWTService
}

func NewController(u Usecases.UserUsecase, t Usecases.TaskUsecase, j Infrastructure.JWTService) *Controller {
	return &Controller{userUC: u, taskUC: t, jwtSvc: j}
}

// --- Auth / User endpoints ---

type registerReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (ctr *Controller) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload", "details": err.Error()})
		return
	}
	user, err := ctr.userUC.Register(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"username": user.Username, "role": user.Role})
}

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (ctr *Controller) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	user, err := ctr.userUC.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	token, err := ctr.jwtSvc.GenerateToken(user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "username": user.Username, "role": user.Role})
}

// Promote (admin only)
func (ctr *Controller) Promote(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username required"})
		return
	}
	if err := ctr.userUC.Promote(username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "promoted", "username": username})
}

// --- Task endpoints ---

type createTaskReq struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	Status      string `json:"status" binding:"required"`
}

func (ctr *Controller) CreateTask(c *gin.Context) {
	var req createTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload", "details": err.Error()})
		return
	}
	task := Domain.Task{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Status:      req.Status,
	}
	created, err := ctr.taskUC.CreateTask(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

func (ctr *Controller) GetTasks(c *gin.Context) {
	tasks, err := ctr.taskUC.ListTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (ctr *Controller) GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	task, err := ctr.taskUC.GetTaskByID(objID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

type updateTaskReq struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	DueDate     *string `json:"due_date"`
	Status      *string `json:"status"`
}

func (ctr *Controller) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req updateTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	patch := make(map[string]interface{})
	if req.Title != nil {
		patch["title"] = *req.Title
	}
	if req.Description != nil {
		patch["description"] = *req.Description
	}
	if req.DueDate != nil {
		patch["due_date"] = *req.DueDate
	}
	if req.Status != nil {
		patch["status"] = *req.Status
	}
	updated, err := ctr.taskUC.UpdateTask(objID, patch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (ctr *Controller) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := ctr.taskUC.DeleteTask(objID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
