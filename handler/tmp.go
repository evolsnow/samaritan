package handler

import (
	"encoding/json"
	"fmt"
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/base"
	"github.com/evolsnow/samaritan/conn"
	"github.com/garyburd/redigo/redis"
	"log"
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
	fmt.Println(base.LRUCache.Get("test2"))
	//fmt.Println(ps.Get("authId"))

}

func Pm2(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//	page := r.URL.Query().Get("page")
	//	limit := r.URL.Query().Get("per_page")
	//	fmt.Fprintf(w, page+limit)
	//fmt.Println(base.LRUCache.Get("test2"))
	fmt.Println(conn.Get("test2"))

}

func SetJwt(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tokenString := base.NewToken(123)
	fmt.Fprint(w, tokenString)
}
func Ab(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c := conn.Pool.Get()
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
	c := conn.Pool.Get()
	defer c.Close()
	//	_, err := conn.Do("SET", "username", "evol")
	//	if err != nil {
	//	}
	//	username, _ := redis.String(conn.Do("GET", "username"))
	//	f := map[string]string{"hello": username}
	//	js, _ := json.Marshal(f)

	//
	//		conn.Send("HMSET", "album:1", "title", "Red", "rating", 5)
	//		conn.Send("HMSET", "album:2", "title", "Earthbound", "rating", 1)
	//		conn.Send("HMSET", "album:3", "title", "Beat", "rating", 4)
	//		conn.Send("LPUSH", "albums", "1")
	//		conn.Send("LPUSH", "albums", "2")
	//		conn.Send("LPUSH", "albums", "3")
	//		conn.Do("HMSET", "user", "foo", 10, "bar", 20)
	//	ms := &MyStruct{}
	//ab := &Album{}
	//
	//	reply, err := redis.Values(conn.Do("HGETALL", "hi"))
	//	if err != nil {
	//		log.Println("get error")
	//	}
	//	redis.ScanStruct(reply, ms)
	//	log.Println(*ms)

	//	album, err := redis.Values(conn.Do("HGETALL", "album:1"))
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
