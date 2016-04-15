package middleware

import (
	"net/http"
	"strings"
)

// CTypeMiddleware set response content-type
func CTypeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	if strings.Contains(r.Header.Get("Accept"), "application/json") || r.Header.Get("Accept") == "*/*" {
		w.Header().Set("Content-Type", "application/json")
	}
	//w.Header().Set("Content-Type", "application/json")
	next(w, r)
}
