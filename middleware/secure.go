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
			return
		} else {
			log.Println("fuck")
			return
		}
	}
}
