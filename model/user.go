package model

import "log"

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

//read user's password
func (u *User) GetPassword() (pwd string) {
	pwd, err := ReadPassword(u.Id)
	if err != nil {
		log.Println("Error get user's password:", err)
		return ""
	}
	return
}

//save a new user
func (u *User) Save() (err error) {
	err = createUser(u)
	if err != nil {
		log.Println("Error save user:", err)
		return
	}
	return
}
