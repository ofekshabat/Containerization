package db

import (
	"context"
	"fmt"
	
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoURI     = "mongodb://localhost:27017"
	databaseName = "Mechola"
)

var client *mongo.Client
var database *mongo.Database
var ContainersCollection *mongo.Collection
var ImagesCollection *mongo.Collection
var PackagesCollection *mongo.Collection

// Ensures the mongod service is running and connects to MongoDB
func Connect() error {
	if !isServiceRunning() {
		startService()
	}

	fmt.Printf("Connecting to DB... ")
	ctx := context.Background()

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}

	err = client.Connect(ctx)
	if err != nil {
		return err
	}

	database = client.Database(databaseName)
	ContainersCollection = database.Collection("containers")
	ImagesCollection = database.Collection("images")
	PackagesCollection = database.Collection("packages")
	fmt.Println()
	return nil
}
