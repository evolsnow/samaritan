package handler

import (
	"encoding/json"
	"net/http"
)

//bind json to user defined struct
func parseRequest(w http.ResponseWriter, r *http.Request, ds interface{}) bool {

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(ds)
	if err != nil {
		e := map[string]string{"error": err.Error()}
		msg, _ := json.Marshal(e)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(msg)
		return false
	}
	return true
}

//write user defined struct to client
func generateResponse(w http.ResponseWriter, r *http.Request, ds interface{}) {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(ds)
	if err != nil {
		e := map[string]string{"error": err.Error()}
		msg, _ := json.Marshal(e)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(msg)
	}
}
