package middleware

import (
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/base"
	"net/http"
	"strings"
)

var log = base.Logger

func Auth(h httprouter.Handle) httprouter.Handle {
	//jwt
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Look for an Authorization header
		if ah := r.Header.Get("Authorization"); ah != "" {
			// Should be a bearer token
			if len(ah) > 6 && strings.ToUpper(ah[0:6]) == "BEARER" {
				// Try to get from LRU cache
				if ele, hit := base.LRUCache.Get(ah[7:]); hit {
					ps.Set("authId", ele.(string))
					log.Debug("got token from LRU")
					goto Success
				} else {
					//If failed, parse token and add token to cache
					err := base.ParseToken(ah[7:], &ps)
					if err == nil {
						go base.LRUCache.Add(ah[7:], ps.Get("authId"))
						goto Success
					} else {
						msg := "Invalid Token"
						base.SetError(w, msg, http.StatusUnauthorized)
						return
					}
				}
			} else {
				msg := "Invaild Authorization Method"
				base.SetError(w, msg, http.StatusUnauthorized)
				return
			}
		} else {
			base.SetError(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
	Success:
		h(w, r, ps)
	}
}
