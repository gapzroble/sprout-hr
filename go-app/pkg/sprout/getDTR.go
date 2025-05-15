package sprout

import (
	"context"
	"time"

	"github.com/gapzroble/sprout-hr/pkg/mongodb"
	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DTR struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Date string
	In   *time.Time
	Out  *time.Time

	TTL int64 // TODO
}

func (d DTR) toLocalTime() DTR {
	if d.In != nil {
		in := d.In.In(pht)
		d.In = &in
	}

	if d.Out != nil {
		out := d.Out.In(pht)
		d.Out = &out
	}

	return d
}

var (
	databaseName          = "sprout-hr"
	collectionName        = "dtr"
	holidayCollectionName = "holidays"
)

func GetDTR(ctx context.Context) *DTR {
	var result DTR

	collection := mongodb.Collection(databaseName, collectionName)

	date := Now().Format("2006-01-02")
	filter := bson.D{bson.E{Key: "date", Value: date}}

	log.WithField("date", date).Println("Finding dtr")

	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.WithError(err).Warn("Error finding dtr")
		return nil
	}

	dtr := result.toLocalTime()

	log.WithFields(log.Fields{
		"in":  dtr.In,
		"out": dtr.Out,
	}).Println("DTR result")

	return &dtr
}

func IsHoliday(ctx context.Context) (string, bool) {
	collection := mongodb.Collection(databaseName, holidayCollectionName)

	date := Now().Format("2006-01-02")
	filter := bson.D{bson.E{Key: "date", Value: date}}

	log.WithField("date", date).Println("Is holiday?")

	result := struct {
		Date string
		Name string
	}{}

	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.WithError(err).Warn("Error checking holiday")
		return "", false
	}

	return result.Name, true
}
