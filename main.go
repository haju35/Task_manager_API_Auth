package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/haju35/Task_Manager-API_Auth/data"
    "github.com/haju35/Task_Manager-API_Auth/router"
)

func main() {
    // Load configuration from environment variables
    mongoURI := os.Getenv("MONGO_URI")
    if mongoURI == "" {
        mongoURI = "mongodb://localhost:27017"
    }
    dbName := os.Getenv("MONGO_DB")
    if dbName == "" {
        dbName = "taskdb"
    }
    collName := os.Getenv("MONGO_COLLECTION")
    if collName == "" {
        collName = "tasks"
    }
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    // Initialize MongoDB
    log.Printf("Connecting to MongoDB: %s (db=%s, collection=%s)\n", mongoURI, dbName, collName)
    if err := data.InitMongo(mongoURI, dbName, collName); err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v\n", err)
    }


    // Setup router
    r := router.SetupRouter()

    // Start server in a goroutine
    srvAddr := ":" + port
    go func() {
        if err := r.Run(srvAddr); err != nil {
            log.Fatalf("Failed to run server: %v\n", err)
        }
    }()
    log.Printf("Server running on %s\n", srvAddr)

    // Wait for termination signals for graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutting down server...")

    // Optional: give some time for graceful shutdown
    timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    _ = timeoutCtx

    log.Println("Server stopped")
}
