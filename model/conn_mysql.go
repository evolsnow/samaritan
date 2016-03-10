package model

import (
	"github.com/evolsnow/samaritan/common/dbms"
	"github.com/evolsnow/samaritan/common/log"
)

func CreateUser(u *User) {
	stmt, err := dbms.DB.Prepare("INSERT INTO user(redis_id, pid, sam_id, alias, name, phone, password, email, avatar, school, depart, grade, class, studentNum) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(u.Id, u.Pid, u.SamId, u.Alias, u.Name, u.Phone, u.Password, u.Email, u.Avatar, u.School, u.Department, u.Grade, u.Class, u.StudentNum)
	if err != nil {
		log.Fatal(err)
	}
}

func Test() {
	var name string
	row := dbms.DB.QueryRow("select name from user where id = ?", 1)
	err := row.Scan(&name)
	if err != nil {
		log.Error(err)
	}
	log.Println(name)
}
