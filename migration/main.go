package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define the Todo struct to represent the data
// structure for each document in the collection.
type Todo struct {
	Todo      string    `bson:"todo"`
	Status    string    `bson:"status"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	DeletedAt time.Time `bson:"deleted_at"`
}

func main() {
	// Replace the connection string and
	// database name with your MongoDB details.
	connectionString := "mongodb://localhost:27017"
	databaseName := "example_for_go"

	// Create a MongoDB client.
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal("Error creating MongoDB client:", err)
	}

	// Create a context with a timeout to be used with the client.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB.
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	// Get a handle to the "todos" collection.
	collection := client.Database(databaseName).Collection("todos")

	// Create indexes if needed for better query performance.
	// For example, if you want to query by "status",
	// you can create an index on it.
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "status", Value: 1}},
		Options: options.Index().SetUnique(false),
	}
	_, err = collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Fatal("Error creating index:", err)
	}

	// Optionally, create additional indexes
	// for other fields as needed.
	fmt.Println("Migration complete!")
}
