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

type Todo struct {
	ID        string    `bson:"_id,omitempty"`
	Todo      string    `bson:"todo"`
	Status    string    `bson:"status"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	DeletedAt time.Time `bson:"deleted_at"`
}

func main() {
	// Replace the connection string and database name with your MongoDB details.
	connectionString := "mongodb://localhost:27017"
	databaseName := "example_for_go"

	// Create a MongoDB client with options.
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	defer func() {
		// Disconnect the client when the work is done.
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal("Error disconnecting from MongoDB:", err)
		}
	}()

	// Ping the MongoDB server to check if the connection was successful.
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Error pinging MongoDB:", err)
	}

	// Get a handle to the "todos" collection.
	collection := client.Database(databaseName).Collection("todos")

	// CRUD operations
	insertedID := CreateTodo(collection, "Example Task", "Pending")
	todos := ReadTodos(collection)
	UpdateTodoStatus(collection, insertedID, "Completed")
	SoftDeleteTodo(collection, insertedID)

	fmt.Println(todos)
	fmt.Println("CRUD operations completed successfully!")
}

// CreateTodo creates a new todo in the collection.
func CreateTodo(collection *mongo.Collection, todo, status string) string {
	newTodo := Todo{
		Todo:      todo,
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: time.Time{},
	}
	insertResult, err := collection.InsertOne(context.Background(), newTodo)
	if err != nil {
		log.Fatal("Error creating todo:", err)
	}
	fmt.Println("New Todo ID:", insertResult.InsertedID)
	return insertResult.InsertedID.(primitive.ObjectID).Hex()
}

// ReadTodos retrieves all todos from the collection.
func ReadTodos(collection *mongo.Collection) []Todo {
	filter := bson.M{} // Empty filter to retrieve all todos.
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal("Error retrieving todos:", err)
	}
	defer cursor.Close(context.Background())

	var todos []Todo
	if err := cursor.All(context.Background(), &todos); err != nil {
		log.Fatal("Error decoding todos:", err)
	}
	fmt.Println("All Todos:")
	for _, todo := range todos {
		fmt.Printf("Todo ID: %s, Todo: %s, Status: %s\n", todo.ID, todo.Todo, todo.Status)
	}
	return todos
}

// UpdateTodoStatus updates the status of a todo.
func UpdateTodoStatus(collection *mongo.Collection, todoID string, newStatus string) {
	objectID, err := primitive.ObjectIDFromHex(todoID)
	if err != nil {
		log.Fatal("Invalid Todo ID:", err)
	}
	filter := bson.M{"_id": objectID}
	updateData := bson.M{"$set": bson.M{"status": newStatus, "updated_at": time.Now()}}
	_, err = collection.UpdateOne(context.Background(), filter, updateData)
	if err != nil {
		log.Fatal("Error updating todo:", err)
	}
}

// SoftDeleteTodo soft deletes a todo.
func SoftDeleteTodo(collection *mongo.Collection, todoID string) {
	objectID, err := primitive.ObjectIDFromHex(todoID)
	if err != nil {
		log.Fatal("Invalid Todo ID:", err)
	}
	filter := bson.M{"_id": objectID}
	updateData := bson.M{"$set": bson.M{"deleted_at": time.Now()}}
	_, err = collection.UpdateOne(context.Background(), filter, updateData)
	if err != nil {
		log.Fatal("Error soft deleting todo:", err)
	}
}
