package model

import "log"

type User struct {
	Id         int    `json:"id,omitempty" redis:"id"`
	Alias      string `json:"alias,omitempty" redis:"alias"` //nick name
	Name       string `json:"name,omitempty" redis:"name"`   //real name
	Phone      string `json:"phone,omitempty" redis:"phone"`
	Password   string `json:"passwd,omitempty" redis:"passwd"`
	Avatar     string `json:"avatar,omitempty" redis:"avatar"` //avatar url
	School     string `json:"school,omitempty" redis:"school"`
	Department string `json:"depart,omitempty" redis:"depart"`
	Grade      int    `json:"grade,omitempty" redis:"grade"`
	Class      string `json:"class,omitempty" redis:"class"`
	StudentNum string `json:"stuNum,omitempty" redis:"stuNum"` //1218404001...
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
func (u *User) Save() int {
	//return user id for jwt token use
	return createUser(u)
}
