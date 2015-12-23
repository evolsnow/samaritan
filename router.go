package main

import (
	"github.com/evolsnow/gosqd/handler"
	md "github.com/evolsnow/gosqd/middleware"
	"github.com/julienschmidt/httprouter"
)

var r = httprouter.New()

func getRouter() *httprouter.Router {

	r.GET("/", handler.ProductList)
	r.GET("/test", handler.Test)
	r.GET("/sync", syncProduct)
	r.POST("/md", md.BasicAuth(handler.Hhhh))
	return r
}
