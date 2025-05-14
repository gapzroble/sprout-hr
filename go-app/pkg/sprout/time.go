package sprout

import (
	"time"

	log "github.com/sirupsen/logrus"

	_ "time/tzdata"
)

var (
	pht *time.Location
)

func init() {
	var err error

	pht, err = time.LoadLocation("Asia/Manila")
	if err != nil {
		pht = time.Local
		log.WithError(err).Warn("Failed to load timezone")
	}
	time.Local = pht
}

func Now() time.Time {
	now := time.Now().In(pht)
	// adjust if now 12am - 6am
	if now.Hour() >= 0 && now.Hour() <= 6 {
		adj := time.Date(now.Year(), now.Month(), now.Day()-1, 23, 59, 59, 0, pht)
		log.WithFields(log.Fields{
			"orig":     now,
			"adjusted": adj,
		}).Println("Now")
		return adj
	}
	return now // no adjustment
}
