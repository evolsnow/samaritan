package main

import (
	"github.com/codegangsta/negroni"
	"github.com/evolsnow/samaritan/conn"
	mw "github.com/evolsnow/samaritan/middleware"
	"log"
	"net"
	"os"
	"strconv"
)

func main() {
	config, err := ParseConfig("config.json")
	if err != nil {
		log.Println("a vailid json config file must exist")
		os.Exit(1)
	}
	//init redis pool
	redisPort := strconv.Itoa(config.RedisPort)
	conn.Pool = conn.NewPool(net.JoinHostPort(config.RedisAddress, redisPort), config.RedisPassword, config.RedisDb)

	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.HandlerFunc(mw.CTypeMiddleware),
	)
	r := newRouter()
	n.UseHandler(r)
	srvPort := strconv.Itoa(config.Port)
	n.Run(net.JoinHostPort(config.Server, srvPort))
}
