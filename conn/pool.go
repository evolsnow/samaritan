package conn

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

var Pool *redis.Pool

func init() {
	Pool = newPool("127.0.0.1:6379", "123456", 3)
}

func newPool(server, password string, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			//			if _, err := c.Do("AUTH", password); err != nil {
			//				c.Close()
			//				return nil, err
			//			}
			if _, err := c.Do("SELECT", db); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
	}
}
