package sprout

import (
	"log"
	"time"

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
		log.Println("Failed to load timezone", err)
	}
}

func Now(adjust ...bool) time.Time {
	now := time.Now().In(pht)
	log.Println("Orig Now", now)
	if len(adjust) > 0 && adjust[0] && now.Hour() < 6 { // adjust if now 12am - 6am
		return time.Date(now.Year(), now.Month(), now.Day()-1, 23, 59, 59, 0, pht)
	}
	return now // no adjustment
}
