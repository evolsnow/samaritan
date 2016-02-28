package model

type Mission struct {
	Id            int       `json:"-" redis:"id"`             //private id
	Pid           string    `json:"id,omitempty" redis:"pid"` //public id
	Name          string    `json:"name,omitempty" redis:"name"`
	createTime    int64     `json:"createTime,omitempty" redis:"createTime"`   //create time timestamp of this mission
	Desc          string    `json:"desc,omitempty" redis:"desc"`               //description for the project
	PublisherId   int       `json:"publisherId,omitempty" redis:"publisherId"` //who published the mission
	ReceiversId   []int     `json:"receiversId,omitempty" redis:"-"`           //user list who accepted the mission
	Completion    float32   `json:"completion,omitempty" redis:"completion"`
	CompletedTime int64     `json:"completedTime,omitempty" redis:"completedTime"`
	Comments      []Comment `json:"comments,omitempty" redis:"-"`
}
type Comment struct {
	Id         int    `json:"-" redis:"id"`
	Pid        string `json:"id,omitempty" redis:"pid"`
	When       int64  `json:"when,omitempty" redis:"when"`
	missionPid string `json:"-" redis:"-"`
	CriticPid  string `json:"uid,omitempty" redis:"criticPid"`
	CriticName string `json:"uName,omitempty" redis:"criticName"`
}

func (cm *Comment) Save() {
	createMissionComment(cm)
}
