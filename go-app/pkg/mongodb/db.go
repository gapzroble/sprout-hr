package mongodb

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

func Connect(ctx context.Context, url string) error {
	log.WithField("url", url).Println("Connecting to mongodb")

	var err error

	clientOptions := options.Client().ApplyURI(url)

	cctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	client, err = mongo.Connect(cctx, clientOptions)
	if err != nil {
		log.WithError(err).Error("Error connecting to mongodb")
		return err
	}

	if err = Ping(cctx); err != nil {
		log.WithError(err).Error("Error pinging mongodb")
		return err
	}

	log.Println("MongoClient connected")

	return nil
}

func Ping(ctx context.Context) error {
	return client.Ping(ctx, nil)
}

func Disconnect(ctx context.Context) {
	if client != nil {
		log.Println("Disconnecting mongodb")
		if err := client.Disconnect(ctx); err == nil {
			log.Println("MongoClient disconnected")
		}
	}
}

func Collection(databaseName, collectionName string) *mongo.Collection {
	return client.Database(databaseName).Collection(collectionName)
}
