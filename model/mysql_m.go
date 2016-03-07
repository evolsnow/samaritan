package model

import (
	"database/sql"
	"github.com/evolsnow/samaritan/common/conn"
	"github.com/evolsnow/samaritan/common/log"
	"time"
)

func init() {
	go func() {
		for {
			//leave time for init db
			time.Sleep(time.Second / 2)
			if conn.DB != nil {
				db = conn.DB
				break
			}
		}
	}()
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

func CreateUser(u *User) {
	stmt, err := db.Prepare("INSERT INTO user(redis_id, pid, sam_id, alias, name, phone, password, email, avatar, school, depart, grade, class, studentNum) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(u.Id, u.Pid, u.SamId, u.Alias, u.Name, u.Phone, u.Password, u.Email, u.Avatar, u.School, u.Department, u.Grade, u.Class, u.StudentNum)
	if err != nil {
		log.Fatal(err)
	}
}
