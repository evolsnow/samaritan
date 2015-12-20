package main

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func Test(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	conn := pool.Get()
	defer conn.Close()
	username, err := redis.String(conn.Do("GETSET", "username", "evol"))
	if err != nil {
		log.Println(err)
	}
	f := map[string]string{"hello": username}
	js, _ := json.Marshal(f)
	w.Write([]byte(js))
}
