package handler

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/evolsnow/httprouter"
	"log"
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
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	token.Claims["userId"] = "123"
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte("mySigningKey"))
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Fprint(w, tokenString)
}
