package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
 ID      primitive.ObjectID    `json:"id"`
 Title       string    `json:"title"`
 Description string    `json:"description"`
 DueDate     string `json:"due_date"`
 Status      string    `json:"status"`
}

