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

// Define the Todo struct to represent the data structure
// for each document in the collection.
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

	// Ping the MongoDB server to check if
	// the connection was successful.
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Error pinging MongoDB:", err)
	}

	// Get a handle to the "todos" collection.
	collection := client.Database(databaseName).Collection("todos")

	// Create indexes if needed for better query performance.
	// For example, if you want to query by "status", you can create an index on it.
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "status", Value: 1}},
		Options: options.Index().SetUnique(false),
	}
	_, err = collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatal("Error creating index:", err)
	}

	fmt.Println("Migration complete!")
}
