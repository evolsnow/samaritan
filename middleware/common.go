package middleware

import (
	"net/http"
	"strings"
)

func CTypeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		w.Header().Set("Content-Type", "application/json")
	}
	//w.Header().Set("Content-Type", "application/json")
	next(w, r)
}
