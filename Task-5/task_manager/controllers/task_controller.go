package controllers

import (
    "errors"
    "net/http"
    "task_manager/data"
    "task_manager/models"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type TaskController struct {
    Service *data.TaskService
}

func NewTaskController(service *data.TaskService) *TaskController {
    return &TaskController{Service: service}
}

func (c *TaskController) GetAll(ctx *gin.Context) {
    tasks, err := c.Service.GetAll()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tasks", "details": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, tasks)
}

func (c *TaskController) GetByID(ctx *gin.Context) {
    id := ctx.Param("id")
    task, err := c.Service.GetByID(id)

    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) || err.Error() == "invalid task ID format" || err.Error() == "task not found" {
            ctx.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
        } else {
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch task", "details": err.Error()})
        }
        return
    }

    ctx.JSON(http.StatusOK, task)
}

func (c *TaskController) Create(ctx *gin.Context) {
    var task models.Task

    if err := ctx.ShouldBindJSON(&task); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON payload or missing required fields", "details": err.Error()})
        return
    }

    result, err := c.Service.Create(task)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not create task", "details": err.Error()})
        return
    }

    // The inserted ID is the ObjectID, which is now the task's ID field
    insertedID, ok := result.InsertedID.(primitive.ObjectID)
    if !ok {
        // Fallback error
        ctx.JSON(http.StatusCreated, gin.H{"message": "task created, but ID retrieval failed"})
        return
    }

    // Retrieve the full task to return to the client
    createdTask, err := c.Service.GetByID(insertedID.Hex())
    if err != nil {
        // Log this as it's an unexpected state
        ctx.JSON(http.StatusCreated, gin.H{"message": "task created, but failed to retrieve full object"})
        return
    }
    ctx.JSON(http.StatusCreated, createdTask)
}

func (c *TaskController) Update(ctx *gin.Context) {
    id := ctx.Param("id")
    var updated models.Task

    if err := ctx.ShouldBindJSON(&updated); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON payload or missing required fields", "details": err.Error()})
        return
    }

    task, err := c.Service.Update(id, updated)
    if err != nil {
        if err.Error() == "task not found" || err.Error() == "invalid task ID format" {
            ctx.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
        } else {
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not update task", "details": err.Error()})
        }
        return
    }

    ctx.JSON(http.StatusOK, task)
}

func (c *TaskController) Delete(ctx *gin.Context) {
    id := ctx.Param("id")

    err := c.Service.Delete(id)
    if err != nil {
        if err.Error() == "task not found" || err.Error() == "invalid task ID format" {
            ctx.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
        } else {
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete task", "details": err.Error()})
        }
        return
    }

    ctx.Status(http.StatusNoContent)
}
