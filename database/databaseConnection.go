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
    MongoDB := "mongodb://admin:password123@localhost:27017/restraunt?authSource=admin"
	fmt.Println(MongoDB)

	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDB))

	if err != nil {
		log.Fatal(err)
	}
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(c)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connected to mongodb")

	return client
}

var Client *mongo.Client =DBInstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("restraunt").Collection(collectionName)
	return collection
}
