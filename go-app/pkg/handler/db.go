package handler

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
)

// ConnectMongoDb take mongodb url and related to connections
func ConnectMongoDb(ctx context.Context, url string) error {
	log.WithField("url", url).Println("Connecting to mongodb")

	var err error

	clientOptions := options.Client().ApplyURI(url)

	cctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err = mongo.Connect(cctx, clientOptions)
	if err != nil {
		log.WithError(err).Error("Error connecting to mongodb")
		return err
	}

	// Check the connection
	if err = client.Ping(cctx, nil); err != nil {
		log.WithError(err).Error("Error pinging mongodb")
		return err
	}

	log.Println("MongoClient connected")

	return nil
}

func DisconnectMongoDb(ctx context.Context) {
	if client != nil {
		log.Println("Disconnecting mongodb")
		if err := client.Disconnect(ctx); err == nil {
			log.Println("MongoClient disconnected")
		}
	}
}
