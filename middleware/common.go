package middleware

import "net/http"

func CTypeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	if accept := r.Header.Get("Accept"); accept == "application/json" {
		w.Header().Set("Content-Type", "application/json")
	}

	next(w, r)
}
