package main

import (
	"github.com/codegangsta/negroni"
	"github.com/garyburd/redigo/redis"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

var pool *redis.Pool

func main() {
	//	Rds, err := redis.Dial("tcp", "127.0.0.1:6379", 2)
	pool = newPool("127.0.0.1:6379")

	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.HandlerFunc(myMiddleware),
	)

	r := httprouter.New()
	r.GET("/", ProductList)
	r.GET("/test", Test)

	n.UseHandler(r)
	n.Run(":8080")
}

func myMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	if accept := r.Header.Get("Accept"); accept == "application/json" {
		w.Header().Set("Content-Type", "application/json")
	}

	next(w, r)
}

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			//              if _, err := c.Do("AUTH", password); err != nil {
			//                  c.Close()
			//                  return nil, err
			//              }
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
