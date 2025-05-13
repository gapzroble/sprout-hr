package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gapzroble/sprout-hr/pkg/sprout"
)

func Endpoints(w http.ResponseWriter, r *http.Request) {
	defer handlePanic()

	link := getLink(r.URL.Query().Get("apikey"))
	links := link.Build(1)

	w.Write([]byte(links))
}

func getLink(apikey string) (link *Link) {
	link = NewLink("sprout", "Sprout")

	if isWeekend() {
		link.AddChild(NewLink("weekend", "Rest Day"))
		return
	}

	if isHoliday() {
		link.AddChild(NewLink("holiday", "Rest Day (Holiday)"))
		return
	}

	if !sprout.CanLogin() {
		link.AddChild(NewLink("cant_login_yet", "Login later"))
		return
	}

	timeIn, timeOut := sprout.GetDTR()
	if timeIn != nil {
		link.AddChild(NewLink("logged_in", fmt.Sprintf("Logged in (%s)", timeIn.Format("03:04pm"))))
	} else {
		link.AddChild(NewLink("login", "Login", fmt.Sprintf("/login?apikey=%s", apikey)))
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

	link.PrependChild(NewLink("logout", "Logout", fmt.Sprintf("/logout?apikey=%s", apikey)))

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

		log.Println(message)
	}
}
