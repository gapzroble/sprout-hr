package sprout

import (
	"time"
)

var (
	pht *time.Location
)

func init() {
	var err error

	pht, err = time.LoadLocation("Asia/Manila")
	if err != nil {
		pht = time.Local
		logger.Error(&logger.LogEntry{
			Message:      "Failed to load timezone",
			ErrorMessage: err.Error(),
		})
	}
}

func Now(adjust ...bool) time.Time {
	now := time.Now().In(pht)
	logger.Info(&logger.LogEntry{
		Message: "Orig Now",
		Keys: map[string]interface{}{
			"now": now,
		},
	})
	if len(adjust) > 0 && adjust[0] && now.Hour() < 6 { // adjust if now 12am - 6am
		return time.Date(now.Year(), now.Month(), now.Day()-1, 23, 59, 59, 0, pht)
	}
	return now // no adjustment
}
