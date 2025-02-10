package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var client *mongo.Client

func ConnectDB() {
	clientOptions := options.Client().ApplyURI("mongodb://hetic:password@localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB: ", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Error pinging MongoDB: ", err)
	}

	fmt.Println("Connected to MongoDB!")
}
