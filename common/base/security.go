package base

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/scrypt"
	"strconv"
)

const (
	JwtKey          = "36861f1530941263e6f61b43743074d8"
	TokenSalt       = "03e23aeb89f13ff4323e641a559db414"
	PrivateChatSalt = "e1b46b79232e42eb4656ee2c810a1d5b"
	UserIdSalt      = "1d143777c383ec8f7c7b7e2a4879ce85"
	TodoIdSalt      = "f7f32e72f01973acc96e5038113f67e4"
	ProjectIdSalt   = "d27023a4f4939d8059b5eed20e86e6be"
	MissionIdSalt   = "d27023a4f4939d8059b5eed20e86e6be"
	CommentIdSalt   = "d27023a4f4939d8059b5eed20e86e6be"
)

func MakeToken(id int) string {
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	token.Claims["userId"] = id
	token.Claims["salt"] = TokenSalt
	// Sign and get the complete encoded token as a string
	tokenString, _ := token.SignedString([]byte(JwtKey))
	return "Bearer " + tokenString
}

func ParseToken(ah string) (uid int, err error) {
	token, err := jwt.Parse(ah, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtKey), nil
	})
	if err == nil && token.Valid {
		userId := token.Claims["userId"].(float64)
		return int(userId), nil
	} else {
		return 0, err
	}
}

//high level secret,use 'scrypt' instead of hash+salt
func EncryptedPassword(pwd string) string {
	salt := fmt.Sprintf("%s@samaritan.tech", pwd)
	dk, _ := scrypt.Key([]byte(pwd), []byte(salt), 16384, 8, 1, 32)
	//return string(dk)

	h := md5.New()
	h.Write(dk)
	h.Write([]byte(salt))
	encrypted := hex.EncodeToString(h.Sum(nil))
	return encrypted
}

func NewPrivateChatId(raw string) string {
	return hashWithSalt(raw, PrivateChatSalt)
}

func HashedUserId(id int) string {
	raw := strconv.Itoa(id)
	return hashWithSalt(raw, UserIdSalt)
}

func HashedTodoId(id int) string {
	raw := strconv.Itoa(id)
	return hashWithSalt(raw, TodoIdSalt)
}

func HashedProjectId(id int) string {
	raw := strconv.Itoa(id)
	return hashWithSalt(raw, ProjectIdSalt)
}

func HashedMissionId(id int) string {
	raw := strconv.Itoa(id)
	return hashWithSalt(raw, MissionIdSalt)
}

func HashedCommentId(id int) string {
	raw := strconv.Itoa(id)
	return hashWithSalt(raw, CommentIdSalt)
}

func hashWithSalt(raw, salt string) string {
	h := md5.New()
	h.Write([]byte(raw))
	h.Write([]byte(salt))
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
