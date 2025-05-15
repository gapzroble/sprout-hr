package sprout

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gapzroble/sprout-hr/pkg/mongodb"
	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
)

func login(token string) (string, error) {
	data := url.Values{}
	data.Set("typeClock", "ClockIn")
	data.Set("Username", os.Getenv("USERNAME"))
	data.Set("Password", os.Getenv("PASSWORD"))
	data.Set("__RequestVerificationToken", token)
	data.Set("X-Requested-With", "XMLHttpRequest")

	req, err := http.NewRequest("POST", clockInUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return "New request faield", err
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

func Login(ctx context.Context, token string) (string, error) {
	message, err := login(token)
	if err != nil {
		return "Login failed", err
	}

	log.Println("Saving DTR..")

	now := Now()
	dtr := bson.M{
		"date":  now.Format("2006-01-02"),
		"in":    &now,
		"login": message,
	}

	collection := mongodb.Collection(databaseName, collectionName)

	result, err := collection.InsertOne(ctx, dtr)
	if err != nil {
		return "Save dtr faield", err
	}
	log.WithField("insertedId", result.InsertedID).Println("result")

	return message, nil
}
