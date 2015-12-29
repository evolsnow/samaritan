package handler

import (
	"fmt"
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/base"
	"net/http"
)

type requestData struct {
	Jjj        int    `json:"cardNo,omitempty"`
	MethodName string `json:"methodName"`
	Inner      nestedJson
}

type nestedJson struct {
	Name string
	Age  int
}

func Hi(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rd := new(requestData)
	if !parseRequest(w, r, rd) {
		return
	}
	generateResponse(w, r, rd)
}

func Pm(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("per_page")
	fmt.Fprintf(w, page+limit)

}

func SetJwt(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tokenString := base.NewToken("123")
	fmt.Fprint(w, tokenString)
}
