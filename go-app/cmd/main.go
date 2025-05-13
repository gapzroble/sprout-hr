package main

import (
   "fmt"
   "log"
   "net/http"
   "runtime"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
   fmt.Fprint(w, "Nginx server as a reverse proxy!\n")
   fmt.Fprintf(w, "Go version: %s\n", runtime.Version())
}

func main() {
   http.HandleFunc("/", indexHandler)
   log.Fatal(http.ListenAndServe(":3000", nil))
}
