package main

import (
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/handler"
	mw "github.com/evolsnow/samaritan/middleware"
)

var r = httprouter.New()

func newRouter() *httprouter.Router {

	r.GET("/", handler.ProductList)
	r.GET("/sync", syncProduct)

	//test
	r.GET("/test", handler.Test)
	r.POST("/hi", mw.Auth(handler.Hi))
	r.POST("/hin", handler.Hi)
	r.GET("/set", handler.SetJwt)
	r.GET("/pm", handler.Pm)
	r.GET("/ab", handler.Ab)
	return r
}
