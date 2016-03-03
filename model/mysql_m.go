package model

import (
	"github.com/evolsnow/samaritan/conn"
)

func Test() {
	var (
		name string
	)
	db := conn.DB
	rows, err := db.Query("select name from user where pid = ?", 123)
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
