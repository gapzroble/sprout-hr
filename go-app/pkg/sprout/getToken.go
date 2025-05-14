package sprout

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	endpoint   string = "https://xlr8clock.hrhub.ph/WebBundy"
	clockInUrl string = "https://xlr8clock.hrhub.ph/WebBundy/ClockIn"
)

func GetRequestVerificationToken() (string, error) {
	log.Println("Getting requestion verification token ..")

	response, err := http.Get(endpoint)
	if err != nil {
		log.Println("Failed to get token", err)
		return "", err
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Failed to read response", err)
		return "", err
	}

	if response.StatusCode > 299 {
		log.Println("Request failed",
			map[string]interface{}{
				"StatusCode":       response.StatusCode,
				"Response Headers": response.Header,
				"Response Body":    string(responseBody),
			},
		)
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
		log.Println("Found token:", token)

		return token, nil
	}

	return "", nil
}
