package base

import (
	"encoding/json"
	"net/http"
)

//set http status and reply error
func SetError(w http.ResponseWriter, desc string, status int) {
	e := map[string]interface{}{"code": status, "error": desc}
	msg, _ := json.Marshal(e)
	w.WriteHeader(status)
	w.Write(msg)
}
