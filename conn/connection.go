package conn

import (
	"database/sql"
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net"
	"strconv"
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

func NewDB(password, host string, port int, database string) *sql.DB {
	server := net.JoinHostPort(host, strconv.Itoa(port))
	db, _ := sql.Open("mysql", fmt.Sprintf("remote:%s@tcp(%s)/%s?autocommit=true", password, server, database))
	if err := db.Ping(); err != nil {
		log.Fatal("failed to connect mysql:", err)
		return nil
	}
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(100)
	return db
}
