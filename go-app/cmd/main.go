package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gapzroble/sprout-hr/pkg/handler"
	_ "github.com/gapzroble/sprout-hr/pkg/sprout"
	log "github.com/sirupsen/logrus"
)

var (
	apikey string
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	apikey = os.Getenv("APIKEY")
	if apikey == "" {
		apikey = time.Now().String()
	}
	log.WithField("value", apikey).Println("apikey")

	h := sha256.New()
	h.Write([]byte(apikey))
	bs := h.Sum(nil)
	apikey = fmt.Sprintf("%x", bs)
}

func main() {
	ctx := context.Background()

	if err := handler.ConnectMongoDb(ctx, os.Getenv("MONGODB_URL")); err != nil {
		log.Panicln("Failed to connect to mongodb", err)
		return
	}
	defer handler.DisconnectMongoDb(ctx)

	http.HandleFunc("/endpoints", authorized(ctx, handler.Endpoints))
	http.HandleFunc("/login", authorized(ctx, handler.Login))
	http.HandleFunc("/logout", authorized(ctx, handler.Logout))

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func authorized(ctx context.Context, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("---------------------------------------------------------------------------------")
		log.WithField("path", r.URL.Path).Println("Request START:")
		defer log.Println("Request END")

		header := r.Header.Get("Authorization")
		if header != apikey {
			for key, value := range r.Header {
				log.WithField(key, value).Println("header")
			}
			log.WithField("value", header).Warn("Wrong/no apikey")
			w.WriteHeader(404)
			// w.Write([]byte("404 page not found"))
			return
		}

		rctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		next.ServeHTTP(w, r.WithContext(rctx))
	})
}
