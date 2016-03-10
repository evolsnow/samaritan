package handler

import (
	"encoding/json"
	"fmt"
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/common/base"
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
	go model.Test()
	c := dbms.Pool.Get()
	defer c.Close()
	username, err := redis.String(c.Do("GETSET", "username", "evol"))
	if err != nil {
		log.Println(err)
	}
	f := map[string]string{"hello": username}
	js, _ := json.Marshal(f)

	w.Write((js))
}
