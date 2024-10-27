// db/database.go

package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectDB() {
	// Updated URI with authSource=admin
	uri := "mongodb://dev:dev@mongo:27017"

	// Set MongoDB client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}

	// Ping the database to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping:", err)
	}

	fmt.Println("Connected to MongoDB!")
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
