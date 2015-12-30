package middleware

import (
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/base"
	"net/http"
)

func Auth(h httprouter.Handle) httprouter.Handle {
	//jwt
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		err := base.ParseToken(r, &ps)
		if err == nil {
			h(w, r, ps)
		} else {
			base.SetError(w, err.Error(), http.StatusUnauthorized)
		}
	}
}
