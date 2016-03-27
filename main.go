package main

import (
	"flag"
	"github.com/evolsnow/negroni"
	"github.com/evolsnow/samaritan/common/dbms"
	"github.com/evolsnow/samaritan/common/log"
	"github.com/evolsnow/samaritan/common/rpc"
	mw "github.com/evolsnow/samaritan/middleware"
	"net"
	"strconv"
)

const CacheDB = "0"

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
	if cfg.RedisS.DB == CacheDB {
		log.Fatal("redis db can not be same as cache db: '0'")
	}

	//init redis  pool
	redisServer := net.JoinHostPort(cfg.RedisS.Address, strconv.Itoa(cfg.RedisS.Port))
	dbms.Pool = dbms.NewPool(redisServer, cfg.RedisS.Password, cfg.RedisS.DB)
	dbms.CachePool = dbms.NewPool(redisServer, cfg.RedisS.Password, cfg.RedisS.DB)

	//init mysql database
	dbms.DB = dbms.NewDB(cfg.MysqlS.Password, cfg.MysqlS.Address, cfg.MysqlS.Port, cfg.MysqlS.DB)

	//init rpc client, domestic+foreign
	rpcServerD := net.JoinHostPort(cfg.RpcSD.Address, strconv.Itoa(cfg.RpcSD.Port))
	rpc.RpcClientD = rpc.NewClientD(rpcServerD)
	rpcServerF := net.JoinHostPort(cfg.RpcSF.Address, strconv.Itoa(cfg.RpcSF.Port))
	rpc.RpcClientF = rpc.NewClientF(rpcServerF)

	//go func() {
	//	title := "帐号注册"
	//	body := "您好，您的注册验证码是：" + base.RandomCode() + "，有效期为5分钟。（请勿直接回复本邮件）"
	//	go rpc.SendMail("gsc1215225@gmail.com", title, body)
	//go rpc.SendMail("lieyan104545@qq.com", title, body)
	//}()

	//init server
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.HandlerFunc(mw.CTypeMiddleware),
	)
	r := newRouter()
	n.UseHandler(r)
	n.Run(net.JoinHostPort(cfg.HttpS.Address, strconv.Itoa(cfg.HttpS.Port)))
}
