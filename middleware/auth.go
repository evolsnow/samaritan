package middleware

import (
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/common/base"
	"github.com/evolsnow/samaritan/common/caches"
	"github.com/evolsnow/samaritan/common/log"
	"net/http"
	"strconv"
	"strings"
)

func Auth(h httprouter.Handle) httprouter.Handle {
	lru := caches.NewLRUCache(100)
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
					uid, err := base.ParseToken(ah[7:])
					if err == nil {
						ps.Set("authId", strconv.Itoa(uid))
						go lru.Add(ah[7:], uid)
						goto Success
					} else {
						msg := "Invalid Token"
						base.UnAuthErr(w, msg)
						return
					}
				}
			} else {
				msg := "Invaild Authorization Method"
				base.UnAuthErr(w, msg)
				return
			}
		} else {
			base.UnAuthErr(w, http.StatusText(http.StatusUnauthorized))
			return
		}
	Success:
		h(w, r, ps)
	}
}
