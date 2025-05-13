package sprout

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	endpoint   string = "https://xlr8clock.hrhub.ph/WebBundy"
	clockInUrl string = "https://xlr8clock.hrhub.ph/WebBundy/ClockIn"
)

func GetRequestVerificationToken() (string, error) {
	response, err := http.Get(endpoint)
	if err != nil {
		logger.Warn(&logger.LogEntry{
			Message:      "Failed to get token",
			ErrorMessage: err.Error(),
		})
		return "", err
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Warn(&logger.LogEntry{
			Message:      "Failed to read response",
			ErrorMessage: err.Error(),
		})
		return "", err
	}

	if response.StatusCode > 299 {
		logger.Warn(&logger.LogEntry{
			Message: "Request failed",
			Keys: map[string]interface{}{
				"StatusCode":       response.StatusCode,
				"Response Headers": response.Header,
				"Response Body":    string(responseBody),
			},
		})
		return "", fmt.Errorf("expecting 2xx response, got %d", response.StatusCode)
	}

	lines := strings.Split(string(responseBody), "\n")

	for _, line := range lines {
		if !strings.Contains(line, "__RequestVerificationToken") {
			continue
		}
		valIndex := strings.Index(line, `value="`)
		if valIndex == -1 {
			return "", errors.New("cannot find value=")
		}
		valPart := line[valIndex+7:]
		quoteIndex := strings.Index(valPart, `"`)
		if quoteIndex == -1 {
			return "", errors.New("cannot find value quote")
		}

		token := valPart[0:quoteIndex]
		logger.InfoStringf("Found token: %s", token)

		return token, nil
	}

	return "", nil
}
