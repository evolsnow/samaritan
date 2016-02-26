package main

import (
	"github.com/evolsnow/negroni"
	"github.com/evolsnow/samaritan/base"
	"github.com/evolsnow/samaritan/conn"
	mw "github.com/evolsnow/samaritan/middleware"
	"log"
	"net"
	"strconv"
)

const LRUCacheSize = 100
const CacheDB = 0

func main() {
	go base.GenerateAvatar("13295121932")
	config, err := ParseConfig("config.json")
	if err != nil {
		log.Fatal("a vailid json config file must exist")
	}
	if config.RedisDB == CacheDB {
		log.Fatal("redis db can not be same as cache db: '0'")
	}

	//init redis database pool
	redisPort := strconv.Itoa(config.RedisPort)
	redisServer := net.JoinHostPort(config.RedisAddress, redisPort)
	if !conn.Ping(redisServer, config.RedisPassword) {
		log.Fatal("connect to redis server failed")
	}
	conn.Pool = conn.NewPool(redisServer, config.RedisPassword, config.RedisDB)
	conn.CachePool = conn.NewPool(redisServer, config.RedisPassword, CacheDB)

	//init LRU cache and simple redis cache
	base.LRUCache = base.NewLRUCache(LRUCacheSize)
	base.Cache = base.NewCache()

	//init server
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
