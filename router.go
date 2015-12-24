package main

import (
	"github.com/evolsnow/gosqd/handler"
	mw "github.com/evolsnow/gosqd/middleware"
	"github.com/julienschmidt/httprouter"
)

var r = httprouter.New()

func getRouter() *httprouter.Router {

	r.GET("/", handler.ProductList)
	r.GET("/sync", syncProduct)

	//test
	r.GET("/test", handler.Test)
	r.POST("/hi", mw.BasicAuth(handler.Hiii))
	r.GET("/pm", handler.Pm)
	return r
}
