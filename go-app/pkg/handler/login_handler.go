package handler

import (
	"net/http"
	"sync"
	"time"

	"github.com/gapzroble/sprout-hr/pkg/sprout"
	log "github.com/sirupsen/logrus"
)

func Login(w http.ResponseWriter, r *http.Request) {
	defer handlePanic()

	if !sprout.CanLogin() {
		w.Write([]byte("Cannot login yet"))
		return
	}

	ctx := r.Context()

	var wg sync.WaitGroup

	var token string
	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		token, err = sprout.GetRequestVerificationToken()
		if err != nil {
			log.WithError(err).Error("Failed to get request verification token")
		}
	}()

	var timeIn *time.Time
	wg.Add(1)
	go func() {
		defer wg.Done()
		if dtr := sprout.GetDTR(ctx); dtr != nil {
			timeIn = dtr.In
		}
	}()

	wg.Wait()

	if timeIn != nil {
		w.Write([]byte("Already logged in"))
		return
	}

	log.Println("Logging in..")

	message, err := sprout.Login(ctx, token)
	if err != nil {
		log.WithError(err).Error("Failed to login")
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(message))

}
