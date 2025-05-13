package main

import (
	"context"
	"fmt"
	"os"

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
		},
	})

	link := getLink(ctx)
	links := link.Build(1)

	logger.Info(&logger.LogEntry{
		Message: "Links",
		Keys: map[string]interface{}{
			"link":  link,
			"links": links,
		},
	})

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       links,
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}, nil
}

func getLink(ctx context.Context) (link *Link) {
	link = NewLink("sprout", "Sprout")

	if isWeekend() {
		link.AddChild(NewLink("weekend", "Rest Day"))
		return
	}

	if isHoliday() {
		link.AddChild(NewLink("holiday", "Rest Day (Holiday)"))
		return
	}

	if !sprout.CanLogin() {
		link.AddChild(NewLink("cant_login_yet", "Login later"))
		return
	}

	timeIn, timeOut := sprout.GetDTR(ctx)
	if timeIn != nil {
		link.AddChild(NewLink("logged_in", fmt.Sprintf("Logged in (%s)", timeIn.Format("03:04pm"))))
	} else {
		link.AddChild(NewLink("login", "Login", "/login"))
	}

	if timeIn == nil {
		link.AddChild(NewLink("no_logout", "Logout"))
		return
	}

	if timeOut != nil {
		link.AddChild(NewLink("logged_out", fmt.Sprintf("Logged out (%s)", timeOut.Format("03:04pm"))))
		return
	}

	if !sprout.CanLogout() {
		link.AddChild(NewLink("cant_logout_yet", "Logout later"))
		return
	}

	link.PrependChild(NewLink("logout", "Logout", "/logout"))

	return
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
