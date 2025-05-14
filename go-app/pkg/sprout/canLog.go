package sprout

import (
	log "github.com/sirupsen/logrus"
)

// CanLogin after 8am
func CanLogin() (result bool) {
	now := Now()
	defer func() {
		log.WithFields(log.Fields{
			"now":    now,
			"hour":   now.Hour(),
			"result": result,
		}).Println("Can login?")
	}()

	result = now.Hour() >= 8
	return
}

// CanLogout after 10:30pm
func CanLogout() (result bool) {
	now := Now()

	defer func() {
		log.WithFields(log.Fields{
			"now":    now,
			"hour":   now.Hour(),
			"result": result,
		}).Println("Can logout?")
	}()

	if now.Hour() < 22 {
		result = false
		return
	}

	result = now.Hour() >= 22 || now.Hour() < 4 // 4am, won't happen
	return
}
