package base

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
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
	e := map[string]interface{}{"code": status, "msg": desc}
	msg, _ := json.Marshal(e)
	w.WriteHeader(status)
	w.Write(msg)
}

func RandomCodeSix() string {
	rand.Seed(time.Now().UTC().UnixNano())
	code := 100000 + rand.Intn(900000)
	return strconv.Itoa(code)
}

func ValidPhone(phone string) bool {
	pattern, _ := regexp.Compile("(13[0-9]|15[01235678]|17[0-9]|18[0-9]|14[57])[0-9]{8}")
	return pattern.MatchString(phone)
}

func ValidMail(mail string) bool {
	idx := strings.LastIndex(mail, "@")
	if idx < 1 || idx == len(mail)-1 {
		return false
	}
	return true
}

//check http bad request error
func BadReqErrHandle(w http.ResponseWriter, desc string) {
	SetError(w, desc, http.StatusBadRequest)
}

//http 403 error
func ForbidErrorHandler(w http.ResponseWriter, desc string) {
	SetError(w, desc, http.StatusForbidden)
}
