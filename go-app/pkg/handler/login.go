package handler

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gapzroble/sprout-hr/pkg/sprout"
)

func Login(w http.ResponseWriter, r *http.Request) {
	defer handlePanic()

	if !sprout.CanLogin() {
		w.Write([]byte("Cannot login yet"))
		return
	}

	var wg sync.WaitGroup

	var token string
	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		token, err = sprout.GetRequestVerificationToken()
		if err != nil {
			log.Println("Failed to get request verification token", err)
		}
	}()

	var timeIn *time.Time
	wg.Add(1)
	go func() {
		defer wg.Done()
		timeIn, _ = sprout.GetDTR(client)
	}()

	wg.Wait()

	if timeIn != nil {
		w.Write([]byte("Already logged in"))
		return
	}

	log.Println("Logging in..")

	message, err := sprout.Login(client, token)
	if err != nil {
		log.Println("Failed to login", err)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(message))

}
