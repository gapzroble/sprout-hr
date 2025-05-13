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
)

func login(token string) (string, error) {
	data := url.Values{}
	data.Set("typeClock", "ClockIn")
	data.Set("Username", os.Getenv("username"))
	data.Set("Password", os.Getenv("password"))
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

func Login(ctx context.Context, token string) (string, error) {
	message, err := login(token)
	if err != nil {
		return "Login failed", err
	}
	// now := Now()
	// dtr := DTR{
	// 	Date: now.Format("2006-01-02"),
	// 	In:   &now,
	// 	TTL:  now.AddDate(0, 2, 0).Unix(),
	// }

	// av, err := dynamodbattribute.MarshalMap(dtr)
	// if err != nil {
	// 	logger.Error(&logger.LogEntry{
	// 		Message:      "Failed to marshalmap dynamodb attribute",
	// 		ErrorMessage: err.Error(),
	// 	})
	// 	return "Marshalmap failed", err
	// }

	// input := &dynamodb.PutItemInput{
	// 	Item:      av,
	// 	TableName: aws.String(os.Getenv("TABLE_NAME")),
	// }

	// sess := session.Must(session.NewSession())
	// svc := dynamodb.New(sess)

	// _, err = svc.PutItemWithContext(ctx, input)
	// if err != nil {
	// 	logger.Warn(&logger.LogEntry{
	// 		Message:      "Failed to save dtr record",
	// 		ErrorMessage: err.Error(),
	// 		Keys: map[string]interface{}{
	// 			"item": dtr,
	// 		},
	// 	})
	// 	return "Save dtr faield", err
	// }

	return message, nil
}
