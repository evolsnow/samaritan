package main

import (
    "fmt"
    "log"
    "net/http"
   // "github.com/go-sql-driver/mysql"
)

func main() {
    http.HandleFunc("/", handler) // each request calls handler
    log.Fatal(http.ListenAndServe("localhost:8005", nil))
}

// handler echoes the Path component of the request URL r.
func handler(w http.ResponseWriter, r *http.Request) {
     fmt.Fprintf(w, r.RemoteAddr)
}
