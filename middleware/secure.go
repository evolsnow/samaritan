package middleware

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func BasicAuth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if accept := r.Header.Get("Accept"); accept == "application/json" {
			h(w, r, ps)
			return
		} else {
			e := map[string]string{"error": "authentication failed"}
			msg, _ := json.Marshal(e)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(msg)
		}
	}
}
