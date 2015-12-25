package handler

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type requestData struct {
	CardNo     int
	MethodName string
	Inner      nestedJson
}

type nestedJson struct {
	Name string
	Age  int
}

func Hiii(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rd := new(requestData)
	if !parseRequest(w, r, rd) {
		return
	}
	generateResponse(w, r, rd)
	//	fmt.Fprintf(w, rd.Inner.Name)
}

func Pm(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("per_page")
	fmt.Fprintf(w, page+limit)

}