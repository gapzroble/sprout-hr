package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gapzroble/sprout-hr/pkg/sprout"
	log "github.com/sirupsen/logrus"
)

func Endpoints(w http.ResponseWriter, r *http.Request) {
	defer handlePanic()

	link := getLink(r.Context(), r.URL.Query().Get("apikey"))
	links := link.Build(1)

	w.Write([]byte(links))
}

func getLink(ctx context.Context, apikey string) (link *Link) {
	link = NewLink("sprout", "Sprout")

	if isWeekend() {
		link.AddChild(NewLink("weekend", "Rest Day"))
		return
	}

	if name, yes := sprout.IsHoliday(ctx, client); yes {
		holiday := "Rest Day (Holiday)"
		if name != "" {
			holiday = name
		}
		link.AddChild(NewLink("holiday", holiday))
		return
	}

	if !sprout.CanLogin() {
		link.AddChild(NewLink("cant_login_yet", "Login later"))
		return
	}

	timeIn, timeOut := sprout.GetDTR(ctx, client)
	if timeIn != nil {
		link.AddChild(NewLink("logged_in", fmt.Sprintf("Logged in (%s)", timeIn.Format("03:04pm"))))
	} else {
		link.AddChild(NewLink("login", "Login", "/login"))
	}

	if timeIn == nil {
		link.AddChild(NewLink("no_logout", "Logout"))
		return
	}

	if timeOut != nil {
		link.AddChild(NewLink("logged_out", fmt.Sprintf("Logged out (%s)", timeOut.Format("03:04pm"))))
		return
	}

	if !sprout.CanLogout() {
		link.AddChild(NewLink("cant_logout_yet", "Logout later"))
		return
	}

	link.PrependChild(NewLink("logout", "Logout", "/logout"))

	return
}

func handlePanic() {
	msg := recover()
	if msg != nil {
		message := "Go panic"
		switch msg := msg.(type) {
		case string:
			message = msg
		case error:
			message = msg.Error()

		default:
			message = "Unknown error type"
		}

		log.Error(message)
	}
}
