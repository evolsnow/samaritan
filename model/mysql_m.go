package model

import (
	"database/sql"
	"github.com/evolsnow/samaritan/conn"
)

func init() {
	for {
		if conn.DB != nil {
			db = conn.DB
			break
		}
	}
}

var db *sql.DB

func Test() {
	var name string
	rows, err := db.Query("select name from user where id = ?", 1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
