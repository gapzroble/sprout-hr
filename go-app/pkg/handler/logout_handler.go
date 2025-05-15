package handler

import (
	"net/http"
	"sync"

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

	var dtr *sprout.DTR

	wg.Add(1)
	go func() {
		defer wg.Done()
		dtr = sprout.GetDTR(ctx)
		if dtr == nil {
			dtr = &sprout.DTR{} // zero value
		}
	}()

	wg.Wait()

	if dtr.Out != nil {
		w.Write([]byte("Already logged out"))
		return
	}

	if dtr.In == nil {
		w.Write([]byte("Not logged in"))
		return
	}

	log.Println("Logging out..")

	message, err := sprout.Logout(ctx, *dtr, token)
	if err != nil {
		log.WithError(err).Error("Failed to logout")
		w.Write([]byte(err.Error()))
	}

	w.Write([]byte(message))
}
