package main

import (
	"github.com/codegangsta/negroni"
	mw "github.com/evolsnow/gosqd/middleware"
)

func main() {
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.HandlerFunc(mw.CTypeMiddleware),
	)
	r := getRouter()
	n.UseHandler(r)
	n.Run(":8080")
}
