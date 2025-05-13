package main

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/gapzroble/mygarminhttpclient/pkg/sprout"
)

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	lctx, _ := lambdacontext.FromContext(ctx)
	logger.Init(lctx.AwsRequestID, os.Getenv("AWS_LAMBDA_FUNCTION_NAME"))

	defer handlePanic()

	logger.Info(&logger.LogEntry{
		Message: "Got event",
		Keys: map[string]interface{}{
			"event": event,
			"dt":    sprout.Now().Format(time.RFC3339),
		},
	})

	if !sprout.CanLogout() {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "Cannot logout yet",
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
		}, nil
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		ssm := ssmhelper.NewWithCacheCheck(ctx)
		if err := ssm.UpdateEnvironment(os.Getenv("SSMPath"), true, true); err != nil {
			logger.Error(&logger.LogEntry{
				Message:      "Failed to get parameters",
				ErrorMessage: err.Error(),
			})
		}
	}()

	var token string
	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		token, err = sprout.GetRequestVerificationToken()
		if err != nil {
			logger.Error(&logger.LogEntry{
				Message:      "Failed to get request verification token",
				ErrorMessage: err.Error(),
			})
		}
	}()

	var timeIn *time.Time
	var timeOut *time.Time
	wg.Add(1)
	go func() {
		defer wg.Done()
		timeIn, timeOut = sprout.GetDTR(ctx)
	}()

	wg.Wait()

	if timeOut != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "Already logged out",
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
		}, nil
	}

	if timeIn == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "Not logged in",
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
		}, nil
	}

	logger.InfoString("Logging out..")

	message, err := sprout.Logout(ctx, timeIn, token)
	if err != nil {
		logger.Error(&logger.LogEntry{
			Message:      "Failed to logout",
			ErrorMessage: err.Error(),
		})
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       err.Error(),
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       message,
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}

func handlePanic() {
	msg := recover()
	if msg != nil {
		entry := &logger.LogEntry{
			Message:   "Go panic",
			ErrorCode: "GoPanic",
		}
		switch msg := msg.(type) {
		case string:
			entry.ErrorMessage = msg
		case error:
			entry.ErrorMessage = msg.Error()

		default:
			entry.ErrorCode = "Unknown error type"
			entry.SetKey("error", msg)
		}

		logger.Error(entry)
	}
}
