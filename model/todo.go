package model

type Todo struct {
	Id        uint64
	StartTime uint64 //start timestamp of this action
	DeadLine  uint64 //end time
	Desc      string //description for the action
	Owner     User   //whose
}
