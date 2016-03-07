package main

import (
	"flag"
	"github.com/evolsnow/negroni"
	"github.com/evolsnow/samaritan/common/caches"
	"github.com/evolsnow/samaritan/common/conn"
	"github.com/evolsnow/samaritan/common/log"
	mw "github.com/evolsnow/samaritan/middleware"
	"net"
	"strconv"
)

const LRUCacheSize = 100
const CacheDB = 0

func main() {
	var debug bool
	var conf string
	flag.BoolVar(&debug, "d", false, "debug mode")
	flag.StringVar(&conf, "c", "config.json", "specify config file")
	flag.Parse()
	//set global log level
	log.SetLogLevel(debug)

	//parse config
	cfg, err := ParseConfig(conf)
	if err != nil {
		log.Fatal("a vailid json config file must exist")
	}
	if cfg.RedisDB == CacheDB {
		log.Fatal("redis db can not be same as cache db: '0'")
	}

	//init redis  pool
	redisPort := strconv.Itoa(cfg.RedisPort)
	redisServer := net.JoinHostPort(cfg.RedisAddress, redisPort)
	conn.Pool = conn.NewPool(redisServer, cfg.RedisPassword, cfg.RedisDB)
	conn.CachePool = conn.NewPool(redisServer, cfg.RedisPassword, CacheDB)

	//init mysql database
	conn.DB = conn.NewDB(cfg.MysqlPassword, cfg.MysqlAddress, cfg.MysqlPort, cfg.MysqlDB)

	//init LRU cache and simple redis cache
	caches.LRUCache = caches.NewLRUCache(LRUCacheSize)
	caches.Cache = caches.NewCache()

	//init server
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.HandlerFunc(mw.CTypeMiddleware),
	)
	r := newRouter()
	n.UseHandler(r)
	srvPort := strconv.Itoa(cfg.Port)
	n.Run(net.JoinHostPort(cfg.Server, srvPort))
}
