package handler

import (
	"net/http"
	"sync"
	"time"

	"github.com/gapzroble/sprout-hr/pkg/sprout"
	log "github.com/sirupsen/logrus"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	defer handlePanic()

	if !sprout.CanLogout() {
		w.Write([]byte("Cannot logout yet"))
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
	var timeOut *time.Time
	wg.Add(1)
	go func() {
		defer wg.Done()
		timeIn, timeOut = sprout.GetDTR(ctx, client)
	}()

	wg.Wait()

	if timeOut != nil {
		w.Write([]byte("Already logged out"))
		return
	}

	if timeIn == nil {
		w.Write([]byte("Not logged in"))
		return
	}

	log.Println("Logging out..")

	message, err := sprout.Logout(ctx, client, timeIn, token)
	if err != nil {
		log.WithError(err).Error("Failed to logout")
		w.Write([]byte(err.Error()))
	}

	w.Write([]byte(message))
}
