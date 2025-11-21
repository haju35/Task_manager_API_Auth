package controllers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/haju35/TaskManager-API/data"
    "github.com/haju35/TaskManager-API/models"
    "go.mongodb.org/mongo-driver/mongo"
)

// CreateTaskHandler handles POST /tasks
func CreateTaskHandler(c *gin.Context) {
    var payload models.Task
    if err := c.ShouldBindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload", "details": err.Error()})
        return
    }

    created, err := data.Create(payload)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task", "details": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, created)
}

// GetAllTasksHandler handles GET /tasks
func GetAllTasksHandler(c *gin.Context) {
    tasks, err := data.GetAll()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tasks", "details": err.Error()})
        return
    }
    c.JSON(http.StatusOK, tasks)
}

// GetTaskHandler handles GET /tasks/:id
func GetTaskHandler(c *gin.Context) {
    id := c.Param("id")

    task, err := data.GetByID(id)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch task", "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, task)
}

// UpdateTaskHandler handles PUT /tasks/:id
func UpdateTaskHandler(c *gin.Context) {
    id := c.Param("id")

    var payload models.Task
    if err := c.ShouldBindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload", "details": err.Error()})
        return
    }

    err := data.Update(id, payload)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task", "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "updated successfully"})
}

// DeleteTaskHandler handles DELETE /tasks/:id
func DeleteTaskHandler(c *gin.Context) {
    id := c.Param("id")

    err := data.Delete(id)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task", "details": err.Error()})
        return
    }

    c.JSON(http.StatusNoContent, gin.H{})
}
