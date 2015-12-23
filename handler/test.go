package handler

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type requestData struct {
	Cardno     int
	Methodname string
}

func Hiii(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rd := new(requestData)
	if !parseRequest(w, r, rd) {
		return
	}
	fmt.Fprintf(w, rd.Methodname)
}
