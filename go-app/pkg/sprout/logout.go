package sprout

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func logout(token string) (string, error) {
	data := url.Values{}
	data.Set("typeClock", "ClockOut")
	data.Set("Username", os.Getenv("username"))
	data.Set("Password", os.Getenv("password"))
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
	log.Println("Response", responseBody)

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

func Logout(client *mongo.Client, timeIn *time.Time, token string) (string, error) {
	message, err := logout(token)
	if err != nil {
		return "Logout failed", err
	}

	now := Now()
	dtr := bson.M{
		"date": now.Format("2006-01-02"),
		"in":   &now,
		"out":  &now,
		"ttl":  now.AddDate(0, 2, 0).Unix(),
	}

	collection := client.Database(databaseName).Collection(collectionName)

	_, err = collection.InsertOne(context.Background(), dtr)
	if err != nil {
		return "Save dtr faield", err
	}

	return message, nil
}
