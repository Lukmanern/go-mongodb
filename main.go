package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx              context.Context
	database         *mongo.Database
	collection       *mongo.Collection
	collectionName   = "todos"
	connectionString = "mongodb://localhost:27017"
	databaseName     = "example_for_go"
)

type Todo struct {
	ID        string    `bson:"_id,omitempty"`
	Todo      string    `bson:"todo"`
	Status    string    `bson:"status"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	DeletedAt time.Time `bson:"deleted_at"`
}

func main() {
	ctx = context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	defer func() {
		// Disconnect the client when the work is done.
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal("Error disconnecting from MongoDB:", err)
		}
	}()

	// Set the MongoDB database and collection
	database = client.Database(databaseName)
	collection = database.Collection(collectionName)

	// CRUD operations
	insertedID := CreateTodo("Example Task", "Pending")
	todos := ReadTodos()
	UpdateTodoStatus(insertedID, "Completed")
	SoftDeleteTodo(insertedID)

	fmt.Println(todos)
	fmt.Println(insertedID)
	fmt.Println("CRUD operations completed successfully!")
}

// GetCollection returns the "todos" collection.
func GetCollection() *mongo.Collection {
	return collection
}

// CreateTodo creates a new todo in the collection.
func CreateTodo(todo, status string) string {
	newTodo := Todo{
		Todo:      todo,
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: time.Time{},
	}
	insertResult, err := collection.InsertOne(ctx, newTodo)
	if err != nil {
		log.Fatal("Error creating todo:", err)
	}
	fmt.Println("New Todo ID:", insertResult.InsertedID)
	return insertResult.InsertedID.(primitive.ObjectID).Hex()
}

// ReadTodos retrieves all todos from the collection.
func ReadTodos() []Todo {
	filter := bson.M{} // Empty filter to retrieve all todos.
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatal("Error retrieving todos:", err)
	}
	defer cursor.Close(ctx)

	var todos []Todo
	if err := cursor.All(ctx, &todos); err != nil {
		log.Fatal("Error decoding todos:", err)
	}
	fmt.Println("All Todos:")
	for _, todo := range todos {
		fmt.Printf("Todo ID: %s, Todo: %s, Status: %s\n", todo.ID, todo.Todo, todo.Status)
	}
	return todos
}

// UpdateTodoStatus updates the status of a todo.
func UpdateTodoStatus(todoID string, newStatus string) {
	objectID, err := primitive.ObjectIDFromHex(todoID)
	if err != nil {
		log.Fatal("Invalid Todo ID:", err)
	}
	filter := bson.M{"_id": objectID}
	updateData := bson.M{"$set": bson.M{"status": newStatus, "updated_at": time.Now()}}
	_, err = collection.UpdateOne(ctx, filter, updateData)
	if err != nil {
		log.Fatal("Error updating todo:", err)
	}
}

// SoftDeleteTodo soft deletes a todo.
func SoftDeleteTodo(todoID string) {
	objectID, err := primitive.ObjectIDFromHex(todoID)
	if err != nil {
		log.Fatal("Invalid Todo ID:", err)
	}
	filter := bson.M{"_id": objectID}
	updateData := bson.M{"$set": bson.M{"deleted_at": time.Now()}}
	_, err = collection.UpdateOne(ctx, filter, updateData)
	if err != nil {
		log.Fatal("Error soft deleting todo:", err)
	}
}
