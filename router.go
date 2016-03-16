package main

import (
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/common/base"
	hd "github.com/evolsnow/samaritan/handler"
	mw "github.com/evolsnow/samaritan/middleware"
	"net/http"
)

//const BaseURL = "/api/1.0"

var r = httprouter.New()

type NotAllowed struct{}
type NotFound struct{}

func newRouter() *httprouter.Router {
	//base url for all requests
	//r.BaseURL = BaseURL

	//user-defined http error handler
	r.MethodNotAllowed = NotAllowed{}
	r.NotFound = NotFound{}

	//http get method

	//http post method
	r.POST("/users", hd.NewUser)
	r.POST("/todos", hd.NewTodo)
	r.POST("/projects", hd.NewProject)
	r.POST("/privateChats", hd.NewPrivateChat)
	//todo upload device token

	//http put method

	//http delete method

	//test
	r.GET("/test", hd.Test)
	r.POST("/hi", hd.Hi)
	r.POST("/hia/:userId", mw.Auth(hd.Hi))
	r.GET("/set", hd.SetJwt)
	r.GET("/pm", hd.Pm)
	r.GET("/pm2", hd.Pm2)
	r.GET("/ab", hd.Ab)

	return r
}

// ServeHTTP makes the NAllowed implement the http.Handler interface.
func (NotAllowed) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	base.SetError(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

// ServeHTTP makes the NFound implement the http.Handler interface.
func (NotFound) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	base.SetError(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}
