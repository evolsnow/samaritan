package main

import (
	"github.com/evolsnow/gosqd/handler"
	mw "github.com/evolsnow/gosqd/middleware"
	"github.com/evolsnow/httprouter"
)

var r = httprouter.New()

func newRouter() *httprouter.Router {

	r.GET("/", handler.ProductList)
	r.GET("/sync", syncProduct)

	//test
	r.GET("/test", handler.Test)
	r.POST("/hi", mw.JwtAuth(handler.Hi))
	r.GET("/set", handler.SetJwt)
	r.GET("/pm", handler.Pm)
	r.GET("/ab", handler.Ab)
	return r
}
