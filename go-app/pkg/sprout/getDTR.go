package sprout

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DTR struct {
	Date string // key
	In   *time.Time
	Out  *time.Time

	TTL int64
}

var (
	databaseName   = "sprout-hr"
	collectionName = "dtr"
)

func GetDTR(client *mongo.Client) (*time.Time, *time.Time) {
	var result DTR

	collection := client.Database(databaseName).Collection(collectionName)

	date := Now().Format("2006-01-02")
	filter := bson.D{bson.E{Key: "date", Value: date}}

	log.Println("Finding dtr", filter)

	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.WithError(err).Warn("Error finding dtr")
		return nil, nil
	}

	if result.In != nil {
		in := result.In.In(pht)
		result.In = &in
	}

	if result.Out != nil {
		out := result.Out.In(pht)
		result.Out = &out
	}

	log.Printf("DTR result (in: %v, out: %v)", result.In, result.Out)

	return result.In, result.Out
}
