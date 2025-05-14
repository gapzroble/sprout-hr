package handler

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
)

// ConnectMongoDb take mongodb url and related to connections
func ConnectMongoDb(url string) error {
	log.Println("Connecting to mongodb", url)

	var err error

	clientOptions := options.Client().ApplyURI(url)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println("Error connecting to mongodb", err)
		return err
	}

	// Check the connection
	if err = client.Ping(ctx, nil); err != nil {
		log.Println("Error pinging mongodb", err)
		return err
	}

	log.Println("MongoClient connected")

	return nil
}

func DisconnectMongoDb() {
	if client != nil {
		log.Println("Disconnecting mongodb")
		if err := client.Disconnect(context.TODO()); err == nil {
			log.Println("MongoClient connected")
		}
	}
}
