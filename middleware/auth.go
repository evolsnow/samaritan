package middleware

import (
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/base"
	"net/http"
)

func LoginRequired(h httprouter.Handle) httprouter.Handle {
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

func Auth(h httprouter.Handle) httprouter.Handle {
	//we need more than "LoginRequired" to access some personal information
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		err := base.ParseToken(r, &ps)
		if err == nil && ps.Get("userId") == ps.Get("mwUserId") {
			h(w, r, ps)
		} else {
			base.SetError(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}
