package handler

import (
	"fmt"
	mw "github.com/evolsnow/gosqd/middleware"
	"github.com/evolsnow/httprouter"

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

func SetJwt(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tokenString := mw.NewToken("123")
	fmt.Fprint(w, tokenString)
}
