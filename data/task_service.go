package data

import (
    "context"
    "errors"
    "time"

    "github.com/haju35/TaskManager-API/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var (
    TasksCollection *mongo.Collection
    mongoClient     *mongo.Client
)

func InitMongo(uri, dbName, collectionName string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    clientOpts := options.Client().ApplyURI(uri)

    client, err := mongo.Connect(ctx, clientOpts)
    if err != nil {
        return err
    }

    if err := client.Ping(ctx, nil); err != nil {
        return err
    }

    mongoClient = client
    TasksCollection = client.Database(dbName).Collection(collectionName)
    return nil
}

// GetAll returns all tasks from MongoDB.
func GetAll() ([]models.Task, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := TasksCollection.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var results []models.Task
    if err := cursor.All(ctx, &results); err != nil {
        return nil, err
    }

    return results, nil
}

// GetByID returns one task by ID.
func GetByID(id string) (*models.Task, error) {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, errors.New("invalid id")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var task models.Task
    err = TasksCollection.FindOne(ctx, bson.M{"_id": oid}).Decode(&task)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, errors.New("not found")
        }
        return nil, err
    }

    return &task, nil
}

// Create inserts a new task.
func Create(t models.Task) (*models.Task, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()


    _, err := TasksCollection.InsertOne(ctx, t)
    if err != nil {
        return nil, err
    }

    return &t, nil
}

// Update modifies an existing task.
func Update(id string, updated models.Task) error {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return errors.New("invalid id")
    }

    updateFields := bson.M{}
    if updated.Title != "" {
        updateFields["title"] = updated.Title
    }
    if updated.Description != "" {
        updateFields["description"] = updated.Description
    }
    if updated.DueDate != "" {
        updateFields["dueDate"] = updated.DueDate
    }
    if updated.Status != "" {
        updateFields["status"] = updated.Status
    }

    if len(updateFields) == 0 {
        return errors.New("no fields to update")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    res, err := TasksCollection.UpdateOne(
        ctx,
        bson.M{"_id": oid},
        bson.M{"$set": updateFields},
    )
    if err != nil {
        return err
    }

    if res.MatchedCount == 0 {
        return errors.New("not found")
    }

    return nil
}

// Delete removes a task by ID.
func Delete(id string) error {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return errors.New("invalid id")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    res, err := TasksCollection.DeleteOne(ctx, bson.M{"_id": oid})
    if err != nil {
        return err
    }

    if res.DeletedCount == 0 {
        return errors.New("not found")
    }

    return nil
}
