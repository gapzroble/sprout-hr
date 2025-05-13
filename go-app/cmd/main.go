package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gapzroble/sprout-hr/pkg/handler"
)

var (
	apikey string
)

func init() {
	apikey = os.Getenv("APIKEY")
	if apikey == "" {
		apikey = time.Now().String()
	}
	log.Println("apikey", apikey)
}

func main() {
	handler.Init(os.Getenv("MONGODB_URL"))

	http.HandleFunc("/endpoints", authorized(handler.Endpoints))
	http.HandleFunc("/login", authorized(handler.Login))
	http.HandleFunc("/logout", authorized(handler.Logout))
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func authorized(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request start: %s\n", r.URL.Path)
		defer log.Printf("Request end: %s\n", r.URL.Path)

		if r.URL.Query().Get("apikey") != apikey {
			response := "Incorrect or no apikey"
			w.WriteHeader(401)
			w.Write([]byte(response))
			return
		}

		next.ServeHTTP(w, r)
	})
}
