package base

import (
	"encoding/json"
	"net/http"
)

//set http status and reply error
func SetError(w http.ResponseWriter, err error, status int) {
	e := map[string]string{"error": err.Error()}
	msg, _ := json.Marshal(e)
	w.WriteHeader(status)
	w.Write(msg)
}
