package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(EnvMONGOURI()))

	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	// ping database
	err = client.Ping(ctx, nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")

	return client
}

// client instance
var DB *mongo.Client = ConnectDB()

// getting database collection
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("socmedAPIDB").Collection(collectionName)

	return collection
}
