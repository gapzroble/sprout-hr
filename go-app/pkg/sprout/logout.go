package sprout

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func logout(token string) (string, error) {
	data := url.Values{}
	data.Set("typeClock", "ClockOut")
	data.Set("Username", os.Getenv("USERNAME"))
	data.Set("Password", os.Getenv("PASSWORD"))
	data.Set("__RequestVerificationToken", token)
	data.Set("X-Requested-With", "XMLHttpRequest")

	req, err := http.NewRequest("POST", clockInUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return "New request failed", err
	}
	req.Header.Set("Content-Type", " application/x-www-form-urlencoded; charset=UTF-8")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "Do http failed", err
	}

	responseBody, _ := io.ReadAll(response.Body)
	log.WithField("body", string(responseBody)).Debug("Response")

	if response.StatusCode > 299 {
		return fmt.Sprintf("%d error", response.StatusCode), fmt.Errorf("expecting 2xx response , got %d", response.StatusCode)
	}

	res, err := NewResponse(responseBody)
	if err != nil {
		return "Read response failed", err
	}

	if !res.Success {
		return "", fmt.Errorf(res.Message)
	}

	return res.Message, nil
}

func Logout(ctx context.Context, client *mongo.Client, dtr *DTR, token string) (string, error) {
	message, err := logout(token)
	if err != nil {
		return "Logout failed", err
	}

	log.Println("Saving DTR..")

	collection := client.Database(databaseName).Collection(collectionName)

	opts := options.Update().SetUpsert(true)
	filter := bson.D{bson.E{Key: "_id", Value: dtr.ID}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: "out", Value: Now()}}}}

	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return "Update dtr failed", err
	}
	log.WithField("result", fmt.Sprintf("%v", result)).Println("result")

	return message, nil
}
