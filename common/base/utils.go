/*
Package base provides some base functions such as random code and valid phoneNumber check
*/
package base

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/evolsnow/samaritan/common/log"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

const (
	AvatarPath = "static/avatar/"
)

//generate avatar url
func getGravatarUrl(phone string) string {
	email := fmt.Sprintf("%s@samaritan.tech", phone)
	h := md5.New()
	h.Write([]byte(email))
	hashed := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("https://cn.gravatar.com/avatar/%s.jpg?d=retro&s=40", hashed)
}

// GenerateAvatar saves avatar file and return relative path
func GenerateAvatar(phone string) (string, error) {
	resp, err := http.Get(getGravatarUrl(phone))
	if err == nil {
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			//use default avatar
			return "", err
		}
		path := fmt.Sprintf("%s%s.jpg", AvatarPath, phone)
		err = ioutil.WriteFile(path, data, 0644)
		if err != nil {
			return "", err
		}
		return path, nil
	}
	return "", err
}

// RandomCodeSix generates random 6-length code
func RandomCodeSix() string {
	rand.Seed(time.Now().UTC().UnixNano())
	code := 100000 + rand.Intn(900000)
	return strconv.Itoa(code)
}

// ValidPhone checks if the phone number is legal
func ValidPhone(phone string) bool {
	pattern := regexp.MustCompile("(13[0-9]|15[01235678]|17[0-9]|18[0-9]|14[57])[0-9]{8}")
	return pattern.MatchString(phone)
}

// ValidMail checks if the mail format is legal
func ValidMail(mail string) bool {
	pattern := regexp.MustCompile("[_a-z0-9-]+(\\.[_a-z0-9-]+)*@[a-z0-9-]+(\\.[a-z0-9-]+)*(\\.[a-z]{2,4})")
	return pattern.MatchString(mail)
}

// ValidSamId checks if the sam id is legal
func ValidSamId(samId string) bool {
	pattern := regexp.MustCompile("^(\\w)+$")
	return pattern.MatchString(samId)
}

// BadReqErr means 400 bad request error
func BadReqErr(w http.ResponseWriter, desc string) {
	setError(w, desc, http.StatusBadRequest)
}

// ForbidErr means 403 forbidden error
func ForbidErr(w http.ResponseWriter, desc string) {
	setError(w, desc, http.StatusForbidden)
}

// MethodNAErr means 405 method not allowed error
func MethodNAErr(w http.ResponseWriter, desc string) {
	setError(w, desc, http.StatusMethodNotAllowed)
}

// UnAuthErr means 401 unauthorized error
func UnAuthErr(w http.ResponseWriter, desc string) {
	setError(w, desc, http.StatusUnauthorized)
}

// NotFoundErr means 404 not found error
func NotFoundErr(w http.ResponseWriter, desc string) {
	setError(w, desc, http.StatusNotFound)
}

// InternalErr means 500 internal error
func InternalErr(w http.ResponseWriter, desc string) {
	setError(w, desc, http.StatusInternalServerError)
}

//set http status and reply error
func setError(w http.ResponseWriter, desc string, status int) {
	e := map[string]interface{}{"code": status, "msg": desc}
	msg, _ := json.Marshal(e)
	log.DebugJson(e)
	w.WriteHeader(status)
	w.Write(msg)
}

// InIntSlice checks if an int is in the slice []int
func InIntSlice(a int, list []int) bool {
	for _, v := range list {
		if a == v {
			return true
		}
	}
	return false
}
