package sprout

// after 8am
func CanLogin() bool {
	return Now(true).Hour() >= 8
}

// after 10:30pm
func CanLogout() bool {
	now := Now(true)

	logger.Info(&logger.LogEntry{
		Message: "Now",
		Keys: map[string]interface{}{
			"now": now,
		},
	})

	if now.Hour() < 22 {
		return false
	}

	return now.Hour() >= 22 || now.Hour() < 4 // 4am, won't happen
}
