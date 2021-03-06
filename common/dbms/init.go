package dbms

import (
	"database/sql"
	"fmt"
	"github.com/evolsnow/samaritan/common/log"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"net"
	"strconv"
	"time"
)

var Pool *redis.Pool

var CachePool *redis.Pool

// NewPool return a redis pool
func NewPool(server, password string, db string) *redis.Pool {
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

var DB *sql.DB

// NewDB return a mysql connection
func NewDB(user, password, host string, port int, database string) *sql.DB {
	server := net.JoinHostPort(host, strconv.Itoa(port))
	db, _ := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?autocommit=true", user, password, server, database))
	if err := db.Ping(); err != nil {
		log.Error("failed to connect mysql:", err)
		return nil
	}
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(100)
	return db
}
