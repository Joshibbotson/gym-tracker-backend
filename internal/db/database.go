package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectDB() {
	// env variable
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI environment variable is not set")
	}
	// uri := "mongodb://localhost:27017/gym-tracker"

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}

	// Ping the database to verify connection
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		panic(err)
	}

	fmt.Println("Connected to MongoDB")
	Client = client
}

func DisconnectDB() {
	if Client != nil {
		err := Client.Disconnect(context.TODO())
		if err != nil {
			log.Fatal("Failed to disconnect:", err)
		}
		fmt.Println("Disconnected from MongoDB.")
	}
}
