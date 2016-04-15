/*
Copyright (c) 2016 samaritan

Licensed under the MIT License (MIT)
Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
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

const CacheDB = "8"

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
		log.Fatal("redis db can not be same as cache db: '8'")
	}

	//init redis  pool
	redisServer := net.JoinHostPort(cfg.RedisS.Address, strconv.Itoa(cfg.RedisS.Port))
	dbms.Pool = dbms.NewPool(redisServer, cfg.RedisS.Password, cfg.RedisS.DB)
	dbms.CachePool = dbms.NewPool(redisServer, cfg.RedisS.Password, CacheDB)

	//init mysql database
	dbms.DB = dbms.NewDB(cfg.MysqlS.Password, cfg.MysqlS.Address, cfg.MysqlS.Port, cfg.MysqlS.DB)

	//init rpc client, domestic+foreign
	rpcServerD := net.JoinHostPort(cfg.RpcSD.Address, strconv.Itoa(cfg.RpcSD.Port))
	rpc.RpcClientD = rpc.NewClientD(rpcServerD)
	rpcServerF := net.JoinHostPort(cfg.RpcSF.Address, strconv.Itoa(cfg.RpcSF.Port))
	rpc.RpcClientF = rpc.NewClientF(rpcServerF)

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
