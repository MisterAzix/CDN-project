package app

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"os"
	"time"
)

var client *mongo.Client

func ConnectDB() {
	var db_username = os.Getenv("DB_USERNAME")
	var db_host = os.Getenv("DB_HOST")
	var db_password = os.Getenv("DB_PASSWORD")
	uri := fmt.Sprintf("mongodb://%s:%s@%s", db_username, db_password, db_host)
	log.Println("Connecting to MongoDB with URI:", uri)
	clientOptions := options.Client().ApplyURI(uri)
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

    clientOptions.SetRegistry(bson.NewRegistry().RegisterDecoder(bson.TypeObjectID, bson.NewObjectIDDecoder()).Build())

	fmt.Println("Connected to MongoDB!")
}
