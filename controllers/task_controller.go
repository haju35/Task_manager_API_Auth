package controllers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/haju35/Task_manager_API_Auth/data"
    "github.com/haju35/Task_manager_API_Auth/models"
    "github.com/haju35/Task_manager_API_Auth/middleware"
    "go.mongodb.org/mongo-driver/mongo"
)

// ======================= Task Handlers =======================

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

// ======================= User Handlers =======================

// RegisterHandler handles POST /register
func RegisterHandler(c *gin.Context) {
    var body struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    u, err := data.CreateUser(body.Username, body.Password)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // hide password
    u.Password = ""
    c.JSON(http.StatusCreated, gin.H{"user": u})
}

// LoginHandler handles POST /login
func LoginHandler(c *gin.Context) {
    var body struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := data.Authenticate(body.Username, body.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    // create JWT token
    token, err := middleware.TokenFromUser(user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}

// PromoteHandler handles PUT /users/:id/promote
func PromoteHandler(c *gin.Context) {
    // only admin can call (router ensures RequireRole)
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
        return
    }

    // promote user to admin
    u, err := data.PromoteToAdmin(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    // hide password
    u.Password = ""
    c.JSON(http.StatusOK, gin.H{"user": u})
}
