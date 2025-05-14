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

	TTL int64 // TODO
}

var (
	databaseName          = "sprout-hr"
	collectionName        = "dtr"
	holidayCollectionName = "holidays"
)

func GetDTR(client *mongo.Client) (*time.Time, *time.Time) {
	var result DTR

	collection := client.Database(databaseName).Collection(collectionName)

	date := Now().Format("2006-01-02")
	filter := bson.D{bson.E{Key: "date", Value: date}}

	log.WithField("date", date).Println("Finding dtr")

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

	log.WithFields(log.Fields{
		"in":  result.In,
		"out": result.Out,
	}).Println("DTR result")

	return result.In, result.Out
}

func IsHoliday(client *mongo.Client) (string, bool) {
	collection := client.Database(databaseName).Collection(holidayCollectionName)

	date := Now().Format("2006-01-02")
	filter := bson.D{bson.E{Key: "date", Value: date}}

	log.WithField("date", date).Println("Finding dtr")

	result := struct {
		Date string
		Name string
	}{}

	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.WithError(err).Warn("Error finding holiday")
		return "", false
	}

	return result.Name, true
}
