package base

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/dgrijalva/jwt-go"
	"github.com/evolsnow/httprouter"
	"net/http"
)

const (
	JwtSignKey   = "36861f1530941263e6f61b43743074d8"
	PasswordSalt = "97096a41d7f2d22348ef9b64fbdfa87a"
	TokenSalt    = "03e23aeb89f13ff4323e641a559db414"
)

func NewToken(id string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	token.Claims["userId"] = id
	token.Claims["salt"] = TokenSalt
	// Sign and get the complete encoded token as a string
	tokenString, _ := token.SignedString([]byte(JwtSignKey))
	return tokenString
}

func ParseToken(r *http.Request, ps *httprouter.Params) (err error) {
	token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtSignKey), nil
	})
	if err == nil && token.Valid {
		userId := token.Claims["userId"].(string)
		(*ps).Set("mwUserId", userId)
		return nil
	} else {
		return err
	}
}

func HashedPassword(pwd string) string {
	h := md5.New()
	h.Write([]byte(pwd))
	h.Write([]byte(PasswordSalt))
	hashed := hex.EncodeToString(h.Sum(nil))
	return hashed
}

//func validSign(XSign, userId string) string {
//	//get key and last visit time from redis
//	appKey, lastVisit := conn.GetSignKey(userId)
//	//parse to compare
//	parts := strings.Split(XSign, ".")
//	current, _ := strconv.Atoi(parts[1])
//	last, _ := strconv.Atoi(lastVisit)
//	if len(parts) != 2 {
//		return "Invalid Sign"
//	}
//	//verify sign
//	h := md5.New()
//	h.Write([]byte(appKey + parts[1]))
//	hash := hex.EncodeToString(h.Sum(nil))
//	if parts[0] == hash {
//		if current > last {
//			//update user sign
//			go conn.UpdateSign(userId, parts[1])
//			return ""
//		} else {
//			return "Replay Attacks"
//		}
//
//	} else {
//		return "Incorrect Sign"
//	}
//}
