package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/base"
	"net/http"
)

func Auth(h httprouter.Handle) httprouter.Handle {
	//jwt+sign
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		//verify sign
		//		XSign := r.Header.Get("X-Sign")
		//		if XSign == "" {
		//			base.SetError(w, "Empty Sign", http.StatusUnauthorized)
		//			return
		//		}
		//parse token
		token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
			return []byte(base.JwtSignKey), nil
		})
		if err == nil && token.Valid {
			userId := token.Claims["userId"].(string)
			ps.Set("user_id", userId)
			h(w, r, ps)
			//			if msg := validSign(XSign, userId); msg == "" {
			//				//save user_id in ps for sharing between middleware or handlers
			//				ps.Set("user_id", userId)
			//				h(w, r, ps)
			//			} else {
			//				base.SetError(w, msg, http.StatusUnauthorized)
			//				return
			//			}
		} else {
			base.SetError(w, err.Error(), http.StatusUnauthorized)
		}
	}
}
