package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoClient *mongo.Client
	MongoDatabase *mongo.Database
)

func ConnectDB()  {
	mongoURI := "mongodb://localhost:27017"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second);
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI);

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatalf("Error the connect in mongoDb: %v", err);
	}

	err = client.Ping(ctx, nil);
	if err != nil {
		log.Fatalf("Error the to ping in mongoDB: %v", err);
	}

	MongoClient = client
	MongoDatabase = client.Database("todolist_fiber")
	fmt.Println("MongoDB connected!")
}

func CloseDB()  {
	if (MongoClient != nil) {
		err := MongoClient.Disconnect(context.Background())
		if err != nil {
			log.Fatalf("Error the to close the MongoDB: %v", err);
		}
		fmt.Println("Conection with MongoDB closed");
	}
}

