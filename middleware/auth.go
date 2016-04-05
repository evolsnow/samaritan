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

var lru = caches.NewLRUCache(100)

func Auth(h httprouter.Handle) httprouter.Handle {
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
						lru.Add(ah[7:], strconv.Itoa(uid))
						goto Success
					} else {
						msg := "token错误，请重新登录"
						base.UnAuthErr(w, msg)
						return
					}
				}
			} else {
				msg := "未知认证方式"
				base.UnAuthErr(w, msg)
				return
			}
		} else {
			base.UnAuthErr(w, "鉴权失败，请登录")
			return
		}
	Success:
		h(w, r, ps)
	}
}
