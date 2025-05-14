package main

import (
	"crypto/sha256"
	"fmt"
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

	h := sha256.New()
	h.Write([]byte(apikey))
	bs := h.Sum(nil)
	apikey = fmt.Sprintf("%x", bs)
}

func main() {
	if err := handler.ConnectMongoDb(os.Getenv("MONGODB_URL")); err != nil {
		log.Panicln("Failed to connect to mongodb", err)
		return
	}

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
			response := "404 page not found"
			w.WriteHeader(404)
			w.Write([]byte(response))
			return
		}

		next.ServeHTTP(w, r)
	})
}
