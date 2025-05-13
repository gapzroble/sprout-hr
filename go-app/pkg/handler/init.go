package handler

import (
	"log"

	"github.com/gapzroble/sprout-hr/pkg/mongodb"
)

func Init(url string) {
	log.Println("Connecting to db", url)
	_, err := mongodb.ConnectMongoDb(url)
	if err != nil {
		log.Panic(err)
	}
}
