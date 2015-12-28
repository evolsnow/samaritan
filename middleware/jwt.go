package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/base"
	"net/http"
)

var SignKeyBytes = []byte("mySigningKey")

func JwtAuth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
			return SignKeyBytes, nil
		})
		if err == nil && token.Valid {
			//save user_id in ps for sharing between middleware or handlers
			ps.Set("user_id", token.Claims["userId"].(string))
			h(w, r, ps)
		} else {
			base.SetError(w, err, http.StatusUnauthorized)
		}
	}
}

func NewToken(id string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	token.Claims["userId"] = id
	// Sign and get the complete encoded token as a string
	tokenString, _ := token.SignedString(SignKeyBytes)
	return tokenString
}
