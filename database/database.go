package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbClient *mongo.Client

func InitDB() {
	mongourl := os.Getenv("MONGOURL")
	clientOptions := options.Client().ApplyURI(mongourl)

	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to mongodb")

	dbClient = client
}

func GetClient() *mongo.Client {
	return dbClient
}
