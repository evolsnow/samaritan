package middleware

import (
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/base"
	"github.com/evolsnow/samaritan/common/caches"
	"github.com/evolsnow/samaritan/common/log"
	"net/http"
	"strings"
)

func Auth(h httprouter.Handle) httprouter.Handle {
	lru := caches.LRUCache
	//jwt
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Look for an Authorization header
		if ah := r.Header.Get("Authorization"); ah != "" {
			// Should be a bearer token
			if len(ah) > 6 && strings.ToUpper(ah[0:6]) == "BEARER" {
				// Try to get from LRU cache
				if ele, hit := lru.Get(ah[7:]); hit {
					ps.Set("authId", ele.(string))
					log.Debug("got token from LRU")
					goto Success
				} else {
					//If failed, parse token and add token to cache
					err := base.ParseToken(ah[7:], &ps)
					if err == nil {
						go lru.Add(ah[7:], ps.Get("authId"))
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
