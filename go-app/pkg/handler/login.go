package handler

import (
	"fmt"
	"net/http"
	"runtime"
)

// "github.com/aws/aws-lambda-go/events"
// "github.com/aws/aws-lambda-go/lambda"
// "github.com/aws/aws-lambda-go/lambdacontext"
// "github.com/gapzroble/sprout-hr/pkg/sprout"

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Nginx server as a reverse proxy!\n")
	fmt.Fprintf(w, "Go version: %s\n", runtime.Version())
}

// func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
// 	lctx, _ := lambdacontext.FromContext(ctx)
// 	logger.Init(lctx.AwsRequestID, os.Getenv("AWS_LAMBDA_FUNCTION_NAME"))

// 	defer handlePanic()

// 	logger.Info(&logger.LogEntry{
// 		Message: "Got event",
// 		Keys: map[string]interface{}{
// 			"event": event,
// 			"dt":    sprout.Now().Format(time.RFC3339),
// 		},
// 	})

// 	if !sprout.CanLogin() {
// 		return events.APIGatewayProxyResponse{
// 			StatusCode: 200,
// 			Body:       "Cannot login yet",
// 			Headers: map[string]string{
// 				"Content-Type": "text/plain",
// 			},
// 		}, nil
// 	}

// 	var wg sync.WaitGroup

// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		ssm := ssmhelper.NewWithCacheCheck(ctx)
// 		if err := ssm.UpdateEnvironment(os.Getenv("SSMPath"), true, true); err != nil {
// 			logger.Error(&logger.LogEntry{
// 				Message:      "Failed to get parameters",
// 				ErrorMessage: err.Error(),
// 			})
// 		}
// 	}()

// 	var token string
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		var err error
// 		token, err = sprout.GetRequestVerificationToken()
// 		if err != nil {
// 			logger.Error(&logger.LogEntry{
// 				Message:      "Failed to get request verification token",
// 				ErrorMessage: err.Error(),
// 			})
// 		}
// 	}()

// 	var timeIn *time.Time
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		timeIn, _ = sprout.GetDTR(ctx)
// 	}()

// 	wg.Wait()

// 	if timeIn != nil {
// 		return events.APIGatewayProxyResponse{
// 			StatusCode: 200,
// 			Body:       "Already logged in",
// 			Headers: map[string]string{
// 				"Content-Type": "text/plain",
// 			},
// 		}, nil
// 	}

// 	logger.InfoString("Logging in..")

// 	message, err := sprout.Login(ctx, token)
// 	if err != nil {
// 		logger.Error(&logger.LogEntry{
// 			Message:      "Failed to login",
// 			ErrorMessage: err.Error(),
// 		})
// 		return events.APIGatewayProxyResponse{
// 			StatusCode: 200,
// 			Body:       err.Error(),
// 			Headers: map[string]string{
// 				"Content-Type": "text/plain",
// 			},
// 		}, nil
// 	}

// 	return events.APIGatewayProxyResponse{
// 		StatusCode: 200,
// 		Body:       message,
// 		Headers: map[string]string{
// 			"Content-Type": "text/plain",
// 		},
// 	}, nil
// }

// func main() {
// 	lambda.Start(handler)
// }

// func handlePanic() {
// 	msg := recover()
// 	if msg != nil {
// 		entry := &logger.LogEntry{
// 			Message:   "Go panic",
// 			ErrorCode: "GoPanic",
// 		}
// 		switch msg := msg.(type) {
// 		case string:
// 			entry.ErrorMessage = msg
// 		case error:
// 			entry.ErrorMessage = msg.Error()

// 		default:
// 			entry.ErrorCode = "Unknown error type"
// 			entry.SetKey("error", msg)
// 		}

// 		logger.Error(entry)
// 	}
// }
