package base

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	avatarPath = "static/avatar/"
)

func getGravatarUrl(phone string) string {
	email := fmt.Sprintf("%s@samaritan.tech", phone)
	h := md5.New()
	h.Write([]byte(email))
	hashed := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("https://cn.gravatar.com/avatar/%s.jpg?d=retro&s=40", hashed)
}

func GenerateAvatar(phone string) (string, error) {
	resp, err := http.Get(getGravatarUrl(phone))
	if err == nil {
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			//use default avatar
			return "", err
		}
		path := fmt.Sprintf("%s%s.jpg", avatarPath, phone)
		go ioutil.WriteFile(path, data, 0644)
		return path, nil
	}
	return "", err
}

//set http status and reply error
func SetError(w http.ResponseWriter, desc string, status int) {
	e := map[string]interface{}{"code": status, "error": desc}
	msg, _ := json.Marshal(e)
	w.WriteHeader(status)
	w.Write(msg)
}

func RandomCode() string {
	rand.Seed(time.Now().UTC().UnixNano())
	code := 100000 + rand.Intn(900000)
	return strconv.Itoa(code)
}

//check http bad request error
func BadReqErrHandle(w http.ResponseWriter, desc string) {
	SetError(w, desc, http.StatusBadRequest)
}

//http 403 error
func ForbidErrorHandler(w http.ResponseWriter) {
	SetError(w, "user match error", http.StatusForbidden)
}
