package middleware

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/evolsnow/httprouter"
	"net/http"
)

var SignKeyBytes = []byte("mySigningKey")

func JwtAuth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
			return SignKeyBytes, nil
		})
		if err == nil && token.Valid {
			//save user_id in ps for sharing from middleware or handler
			ps.Set("user_id", token.Claims["userId"].(string))
			h(w, r, ps)
		} else {
			e := map[string]string{"error": "authentication failed"}
			msg, _ := json.Marshal(e)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(msg)
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
