package model

import "reflect"

type User struct {
	Id         int
	Name       string
	Phone      string
	Avatar     string //avatar url
	School     string
	Department string
	Grade      int
	Class      string
	StudentNum string //1218404001...
}

func (u *User) Update(nu User) {
	value := reflect.ValueOf(nu)
	for i := 0; i < value.NumField(); i++ {
		value.Field(i)
	}
}
