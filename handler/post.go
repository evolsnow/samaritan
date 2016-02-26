package handler

import (
	"errors"
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/base"
	"github.com/mholt/binding"
	"log"
	"net/http"
)

var (
	ErrInvalidRequest = errors.New("invalid request data")
)

func NewTodo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Println("calling func: NewTodo")
	pt := new(postTodoRequest)
	errs := binding.Bind(r, pt)
	if errs != nil {
		base.ErrorHandle(w, r, errs)
		return
	}
	resp := new(postTodoResponse)
	go pt.SaveTodo()
	generateResponse(w, r, resp)
}
