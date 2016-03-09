package handler

import (
	"encoding/json"
	"fmt"
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/base"
	"github.com/evolsnow/samaritan/common/caches"
	"github.com/evolsnow/samaritan/common/dbms"
	"github.com/evolsnow/samaritan/common/log"
	"github.com/evolsnow/samaritan/model"
	"github.com/garyburd/redigo/redis"
	"net/http"
)

type requestData struct {
	Jjj        int    `json:"cardNo,omitempty"`
	MethodName string `json:"methodName"`
	Inner      nestedJson
}

type nestedJson struct {
	Name string
	Age  int
}

func Hi(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rd := new(requestData)
	if !parseReq(w, r, rd) {
		return
	}
	makeResp(w, r, rd)
}

func Pm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//	page := r.URL.Query().Get("page")
	//	limit := r.URL.Query().Get("per_page")
	//	fmt.Fprintf(w, page+limit)
	fmt.Println(caches.LRUCache.Get("test2"))
	//fmt.Println(ps.Get("authId"))

}

func Pm2(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//	page := r.URL.Query().Get("page")
	//	limit := r.URL.Query().Get("per_page")
	//	fmt.Fprintf(w, page+limit)
	//fmt.Println(base.LRUCache.Get("test2"))
	fmt.Println(dbms.Get("test2"))

}

func SetJwt(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tokenString := base.NewToken(123)
	fmt.Fprint(w, tokenString)
}
func Ab(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c := dbms.Pool.Get()
	defer c.Close()
	ret, _ := redis.Bytes(c.Do("GET", ":1:product_list"))
	w.Write(ret)
}

//type MyStruct struct {
//	Foo int `redis:"foo"`
//	Bar int `redis:"bar"`
//}
//
//type Album struct {
//	Title  string `redis:"title"`
//	Rating int    `redis:"rating"`
//}

func Test(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	model.Test()
	c := dbms.Pool.Get()
	defer c.Close()
	//	_, err := dbms.Do("SET", "username", "evol")
	//	if err != nil {
	//	}
	//	username, _ := redis.String(dbms.Do("GET", "username"))
	//	f := map[string]string{"hello": username}
	//	js, _ := json.Marshal(f)

	//
	//		dbms.Send("HMSET", "album:1", "title", "Red", "rating", 5)
	//		dbms.Send("HMSET", "album:2", "title", "Earthbound", "rating", 1)
	//		dbms.Send("HMSET", "album:3", "title", "Beat", "rating", 4)
	//		dbms.Send("LPUSH", "albums", "1")
	//		dbms.Send("LPUSH", "albums", "2")
	//		dbms.Send("LPUSH", "albums", "3")
	//		dbms.Do("HMSET", "user", "foo", 10, "bar", 20)
	//	ms := &MyStruct{}
	//ab := &Album{}
	//
	//	reply, err := redis.Values(dbms.Do("HGETALL", "hi"))
	//	if err != nil {
	//		log.Println("get error")
	//	}
	//	redis.ScanStruct(reply, ms)
	//	log.Println(*ms)

	//	album, err := redis.Values(dbms.Do("HGETALL", "album:1"))
	//	if err != nil {
	//		// handle error
	//		log.Println(err)
	//	}

	//	redis.ScanStruct(album, ab)
	//	//log.Println(*ab)
	//
	//	ret, _ := json.Marshal(ab)
	//	w.Write([]byte(ret))

	username, err := redis.String(c.Do("GETSET", "username", "evol"))
	if err != nil {
		log.Println(err)
	}
	f := map[string]string{"hello": username}
	js, _ := json.Marshal(f)

	w.Write((js))
}
