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
}

func Now(adjust ...bool) time.Time {
	now := time.Now().In(pht)
	if len(adjust) > 0 && adjust[0] && now.Hour() < 6 { // adjust if now 12am - 6am
		adj := time.Date(now.Year(), now.Month(), now.Day()-1, 23, 59, 59, 0, pht)
		log.Println("Orig Now", now, "Adjusted Now", adj)
		return adj
	}
	return now // no adjustment
}
