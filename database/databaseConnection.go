package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBInstance() *mongo.Client {
	MongoDB := "mongodb://mongo:27017"
	fmt.Println(MongoDB)

	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDB))

	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connected to mongodb")

	return client
}

var Client *mongo.Client

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection{
	var collection *mongo.Collection =client.Database("restraunt").Collection(collectionName)
	return collection
}