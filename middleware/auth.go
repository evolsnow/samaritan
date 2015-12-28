package middleware

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/dgrijalva/jwt-go"
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/base"
	"github.com/evolsnow/samaritan/conn"
	"net/http"
	"strconv"
	"strings"
)

var SignKeyBytes = []byte("mySigningKey")

func Auth(h httprouter.Handle) httprouter.Handle {
	//jwt+sign
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		//verify sign
		XSign := r.Header.Get("X-Sign")
		if XSign == "" {
			base.SetError(w, "Empty Sign", http.StatusUnauthorized)
			return
		}
		//parse token
		token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
			return SignKeyBytes, nil
		})
		if err == nil && token.Valid {
			userId := token.Claims["userId"].(string)
			if msg := validSign(XSign, userId); msg == "" {
				//save user_id in ps for sharing between middleware or handlers
				ps.Set("user_id", userId)
				h(w, r, ps)
			} else {
				base.SetError(w, msg, http.StatusUnauthorized)
				return
			}
		} else {
			base.SetError(w, err.Error(), http.StatusUnauthorized)
		}
	}
}

func validSign(XSign, userId string) string {
	//get key and last visit time from redis
	appKey, lastVisit := conn.GetSignKey(userId)
	//parse to compare
	parts := strings.Split(XSign, ".")
	current, _ := strconv.Atoi(parts[1])
	last, _ := strconv.Atoi(lastVisit)
	if len(parts) != 2 {
		return "Invalid Sign"
	}
	//verify sign
	h := md5.New()
	h.Write([]byte(appKey + parts[1]))
	hash := hex.EncodeToString(h.Sum(nil))
	if parts[0] == hash {
		if current == last {
			//update user sign
			go conn.UpdateSign(userId, parts[1])
			return ""
		} else {
			return "Replay Attack"
		}

	} else {
		return "Incorrect Sign"
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
