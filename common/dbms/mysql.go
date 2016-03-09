package dbms

import (
	"github.com/evolsnow/samaritan/common/log"
)

var db = DB

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
