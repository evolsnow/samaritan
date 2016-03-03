package conn

import (
	"database/sql"
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

var Pool *redis.Pool

var CachePool *redis.Pool

var DB *sql.DB

func NewPool(server, password string, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if _, err := c.Do("SELECT", db); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
	}
}

func NewDB(password, server string, port int, database string) *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s:@tcp(%s:%d)/%s?autocommit=true", password, server, port, database))
	if err != nil {
		log.Fatal("failed to connect mysql:", err)
		return nil
	}
	return db
}
