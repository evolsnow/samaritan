package model

import "reflect"

type User struct {
	Id         int    `json:"id" redis:"id"`
	Alias      string `json:"alias" redis:"alias"` //nick name
	Name       string `json:"name" redis:"name"`   //real name
	Phone      string `json:"phone" redis:"phone"`
	Password   string `json:"passwd" redis:"passwd"`
	Avatar     string `json:"avatar" redis:"avatar"` //avatar url
	School     string `json:"school" redis:"school"`
	Department string `json:"depart" redis:"depart"`
	Grade      int    `json:"grade" redis:"grade"`
	Class      string `json:"class" redis:"class"`
	StudentNum string `json:"stuNum" redis:"stuNum"` //1218404001...
}

func (u *User) Update(nu User) {
	value := reflect.ValueOf(nu)
	for i := 0; i < value.NumField(); i++ {
		value.Field(i)
	}
}
