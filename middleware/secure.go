package middleware

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func BasicAuth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if accept := r.Header.Get("Accept"); accept == "application/json" {
			log.Println("authenticated")
			h(w, r, ps)
			return
		} else {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}
