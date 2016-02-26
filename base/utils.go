package base

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	avatarPath = "./avatar/"
)

func getGravatarUrl(phone string) string {
	email := fmt.Sprintf("%s@samaritan.tech", phone)
	h := md5.New()
	h.Write([]byte(email))
	hashed := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("https://cn.gravatar.com/avatar/%s.jpg?d=retro&s=40", hashed)
}

func GenerateAvatar(phone string) {
	resp, err := http.Get(getGravatarUrl(phone))
	if err == nil {
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			//use default avatar
			return
		}
		ioutil.WriteFile(fmt.Sprintf("%s%s.jpg", avatarPath, phone), data, 0644)
	}
}

//set http status and reply error
func SetError(w http.ResponseWriter, desc string, status int) {
	e := map[string]interface{}{"code": status, "error": desc}
	msg, _ := json.Marshal(e)
	w.WriteHeader(status)
	w.Write(msg)
}

//check http bad request error
func BadReqErrorHandle(w http.ResponseWriter, desc string) {
	SetError(w, desc, http.StatusBadRequest)
}

//http 403 error
func ForbidErrorHandler(w http.ResponseWriter) {
	SetError(w, "Authorization error", http.StatusBadRequest)
}
