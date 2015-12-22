package main

import (
	"github.com/codegangsta/negroni"
	"net/http"
)

func main() {
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.HandlerFunc(myMiddleware),
	)
	r := getRouter()
	n.UseHandler(r)
	n.Run(":8080")
}

func myMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	if accept := r.Header.Get("Accept"); accept == "application/json" {
		w.Header().Set("Content-Type", "application/json")
	}

	next(w, r)
}
