package handler

import (
	"encoding/json"
	"github.com/evolsnow/samaritan/base"
	"net/http"
)

//bind json to user defined struct
func parseRequest(w http.ResponseWriter, r *http.Request, ds interface{}) bool {

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(ds)
	if err != nil {
		base.SetError(w, err, http.StatusBadRequest)
		return false
	}
	return true
}

//write user defined struct to client
func generateResponse(w http.ResponseWriter, r *http.Request, ds interface{}) {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(ds)
	if err != nil {
		base.SetError(w, err, http.StatusInternalServerError)
	}
}
