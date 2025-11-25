package router

import (
    "github.com/gin-gonic/gin"
    "github.com/haju35/Task_manager_API_Auth/controllers"
    "github.com/haju35/Task_manager_API_Auth/middleware"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    //Public routes
    r.POST("/register", controllers.RegisterHandler)
    r.POST("/login", controllers.LoginHandler)

    //Protected routes
    authorized := r.Group("/")
    authorized.Use(middleware.AuthMiddleware())

    tasks := authorized.Group("/tasks")
    {
        tasks.POST("", controllers.CreateTaskHandler)
        tasks.GET("", controllers.GetAllTasksHandler)
        tasks.GET("/:id", controllers.GetTaskHandler)
        tasks.PUT("/:id", controllers.UpdateTaskHandler)
        tasks.DELETE("/:id", controllers.DeleteTaskHandler)
    }

    // Promote user (admin only)
    authorized.PUT("/users/:id/promote", middleware.RequireRole("admin"), controllers.PromoteHandler)

    return r
}
