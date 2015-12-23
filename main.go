package main

import (
	"github.com/codegangsta/negroni"
	md "github.com/evolsnow/gosqd/middleware"
)

func main() {
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.HandlerFunc(md.CTypeMiddleware),
	)
	r := getRouter()
	n.UseHandler(r)
	n.Run(":8080")
}
