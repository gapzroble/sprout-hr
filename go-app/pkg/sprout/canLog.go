package sprout

import (
	"log"
)

// after 8am
func CanLogin() bool {
	return Now(true).Hour() >= 8
}

// after 10:30pm
func CanLogout() bool {
	now := Now(true)

	log.Println("Now", now)

	if now.Hour() < 22 {
		return false
	}

	return now.Hour() >= 22 || now.Hour() < 4 // 4am, won't happen
}
