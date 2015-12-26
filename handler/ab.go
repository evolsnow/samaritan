package handler

import (
	"github.com/evolsnow/gosqd/conn"
	"github.com/evolsnow/httprouter"
	"github.com/garyburd/redigo/redis"
	"net/http"
)

func Ab(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c := conn.Pool.Get()
	defer c.Close()
	ret, _ := redis.Bytes(c.Do("GET", ":1:product_list"))
	w.Write(ret)
}
