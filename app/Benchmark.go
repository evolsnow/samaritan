package main

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Test(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", "username", "evol")
	if err != nil {
	}
	username, _ := redis.String(conn.Do("GET", "username"))
	f := map[string]string{"hello": username}
	js, _ := json.Marshal(f)
	w.Write([]byte(js))
}
