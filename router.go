package main

import (
	"github.com/evolsnow/gosqd/handler"
	"github.com/julienschmidt/httprouter"
)

var r = httprouter.New()

func getRouter() *httprouter.Router {

	r.GET("/", handler.ProductList)
	r.GET("/test", handler.Test)
	r.GET("/sync", syncProduct)

	return r
}
