package sprout

import (
	"time"
)

type DTR struct {
	Date string // key
	In   *time.Time
	Out  *time.Time

	TTL int64
}

func GetDTR() (*time.Time, *time.Time) {
	return nil, nil
	// date := Now().Format("2006-01-02")
	// input := &dynamodb.GetItemInput{
	// 	TableName: aws.String(os.Getenv("TABLE_NAME")),
	// 	Key: map[string]*dynamodb.AttributeValue{
	// 		"Date": {
	// 			S: aws.String(date),
	// 		},
	// 	},
	// 	ConsistentRead: aws.Bool(true),
	// }
	// logger.InfoString("Find existing dtr")

	// sess := session.Must(session.NewSession())
	// svc := dynamodb.New(sess)

	// result, err := svc.GetItemWithContext(ctx, input)
	// if err != nil {
	// 	logger.Error(&logger.LogEntry{
	// 		Message:      "Failed to get dtr",
	// 		ErrorMessage: err.Error(),
	// 	})
	// 	return nil, nil
	// }

	// var rec DTR
	// if err := dynamodbattribute.UnmarshalMap(result.Item, &rec); err != nil {
	// 	logger.Error(&logger.LogEntry{
	// 		Message:      "Failed to unmarshalmap",
	// 		ErrorMessage: err.Error(),
	// 	})

	// 	return nil, nil
	// }

	// return rec.In, rec.Out
}
